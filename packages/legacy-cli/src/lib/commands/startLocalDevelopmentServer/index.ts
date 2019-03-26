import * as chokidar from 'chokidar'
import * as cosmiconfig from 'cosmiconfig'
import * as getPort from 'get-port'
import * as Logger from 'koa-logger'
import * as Router from 'koa-router'

import { transpileFunctions } from '../../buildArtifact'
import LocationProvider from '../../locationProvider'
import { Config } from '../../types'

import auth from './auth'
import server from './server'
import Storage from './storage'
import { loadUserDefinedData } from './utils'

import debug from '../../logger'
const logger = debug.extend('start')

function requireUncached(module) {
  delete require.cache[require.resolve(module)]
  return require(module)
}

const DEFAULT_PORT = 3000

export default function startLocalDevelopmentServer(
  emitter,
  config: Config,
  locator: LocationProvider,
  logs: boolean = true
) {
  const rootLevel = locator.integrationRoot

  const LOCAL_DEV_CONFIGURATION = 'dev'
  const explorer = cosmiconfig(LOCAL_DEV_CONFIGURATION, {
    searchPlaces: [`config.${LOCAL_DEV_CONFIGURATION}.js`]
  })
  const router = new Router({ prefix: '/api/' })

  return new Promise<string>(async (resolve, reject) => {
    try {
      const { config: devFunctionsContext = {} } = (await explorer.search(rootLevel)) || {}
      const distPath = locator.buildFunctionsResourcePath('dist')

      // tslint:disable-next-line:no-inner-declarations
      async function refreshFunctions() {
        try {
          emitter.emit('start:localServer:generatingFunctions:start')
          await transpileFunctions(locator.srcFunctionsDir, distPath)
          emitter.emit('start:localServer:generatingFunctions:stop')
        } catch (error) {
          emitter.emit('start:localServer:generatingFunctions:failed', { error })
        }
      }
      await refreshFunctions()

      chokidar
        .watch('.', {
          ignored: /(^|[\/\\])\../,
          cwd: locator.srcFunctionsDir,
          ignoreInitial: true,
          persistent: true,
          followSymlinks: false
        })
        .on('add', refreshFunctions)
        .on('change', refreshFunctions)

      const storage = Storage()
      const expectedPort = process.env.PORT || DEFAULT_PORT

      const port = await getPort({ port: expectedPort })

      // tslint:disable-next-line:no-http-string
      const host = `http://localhost:${port}`

      if (expectedPort !== DEFAULT_PORT) {
        emitter.emit('start:localServer:customPort', { port, host })
      }

      if (port !== Number(expectedPort)) {
        throw `Could not start local server port ${expectedPort} is already in use. You can specify your own port by running PORT=3322 yarn bearer start`
      }
      const bearerBaseURL = `${host}/`
      process.env.bearerBaseURL = bearerBaseURL

      router.post(
        `v3/functions/${config.integrationUuid}/:functionName`,
        functionHandler(distPath, devFunctionsContext, bearerBaseURL),
        (ctx, _next) => ctx.ok(ctx.funcDatum)
      )

      router.post(
        `v2/functions/${config.integrationUuid}/:functionName`,
        functionHandler(distPath, devFunctionsContext, bearerBaseURL),
        (ctx, _next) => ctx.ok(ctx.funcDatum)
      )

      router.all(
        `v1/${config.integrationUuid}/:functionName`,
        functionHandler(distPath, devFunctionsContext, bearerBaseURL),
        (ctx, _next) => ctx.ok(ctx.funcDatum)
      )

      if (logs) {
        server.use(Logger())
      }
      server.use(storage.routes())
      server.use(storage.allowedMethods())
      server.use(auth.routes())
      server.use(auth.allowedMethods())
      server.use(router.routes())
      server.use(router.allowedMethods())

      server.listen(port, () => {
        emitter.emit('start:localServer:start', { port })
        if (/bearer:/.test(process.env.DEBUG || '')) {
          emitter.emit('start:localServer:endpoints', {
            endpoints: [...storage.stack, ...auth.stack, ...router.stack]
          })
        }
      })

      resolve(bearerBaseURL)
    } catch (e) {
      reject(e)
    }
  })
}

const functionHandler = (distPath: string, devFunctionsContext, bearerBaseURL: string) => async (ctx, next) =>
  new Promise(async (resolve, _reject) => {
    try {
      const func = requireUncached(`${distPath}/${ctx.params.functionName}`).default

      const userDefinedData = await loadUserDefinedData({ query: ctx.query })

      const datum = await func.init()(
        {
          context: {
            ...devFunctionsContext.global,
            ...devFunctionsContext[ctx.params.functionName],
            bearerBaseURL,
            ...userDefinedData
          },
          queryStringParameters: ctx.query,
          body: JSON.stringify(ctx.request.body)
        },
        {}
      )
      ctx.funcDatum = datum
      next()
      resolve()
    } catch (e) {
      logger('ERROR: %j', e)
      if (e.code === 'MODULE_NOT_FOUND') {
        ctx.funcDatum = { error: `Function '${ctx.params.functionName}' Not Found` }
      } else {
        ctx.funcDatum = { error: e }
      }
      await next()
      resolve()
    }
  })
