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
      const port = await getPort({ port: 3000 })
      // tslint:disable-next-line:no-http-string
      const bearerBaseURL = `http://localhost:${port}/`
      process.env.bearerBaseURL = bearerBaseURL

      router.post(
        `v3/intents/${config.integrationUuid}/:intentName`,
        intentHandler(distPath, devFunctionsContext, bearerBaseURL),
        (ctx, _next) => ctx.ok(ctx.intentDatum)
      )

      router.post(
        `v2/intents/${config.integrationUuid}/:intentName`,
        intentHandler(distPath, devFunctionsContext, bearerBaseURL),
        (ctx, _next) => ctx.ok(ctx.intentDatum)
      )

      router.all(
        `v1/${config.integrationUuid}/:intentName`,
        intentHandler(distPath, devFunctionsContext, bearerBaseURL),
        (ctx, _next) => ctx.ok(ctx.intentDatum)
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
        emitter.emit('start:localServer:endpoints', {
          endpoints: [...storage.stack, ...auth.stack, ...router.stack]
        })
      })

      resolve(bearerBaseURL)
    } catch (e) {
      reject(e)
    }
  })
}

const intentHandler = (distPath: string, devFunctionsContext, bearerBaseURL: string) => async (ctx, next) =>
  new Promise(async (resolve, _reject) => {
    try {
      const func = requireUncached(`${distPath}/${ctx.params.intentName}`).default

      const userDefinedData = await loadUserDefinedData({ query: ctx.query })

      const datum = await intent.init()(
        {
          context: {
            ...devFunctionsContext.global,
            ...devFunctionsContext[ctx.params.intentName],
            bearerBaseURL,
            ...userDefinedData
          },
          queryStringParameters: ctx.query,
          body: JSON.stringify(ctx.request.body)
        },
        {}
      )
      ctx.intentDatum = datum
      next()
      resolve()
    } catch (e) {
      logger('ERROR: %j', e)
      if (e.code === 'MODULE_NOT_FOUND') {
        ctx.intentDatum = { error: `Function '${ctx.params.intentName}' Not Found` }
      } else {
        ctx.intentDatum = { error: e }
      }
      await next()
      resolve()
    }
  })
