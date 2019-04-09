import * as chokidar from 'chokidar'
import getPort from 'get-port'
import * as Logger from 'koa-logger'
import * as Router from 'koa-router'
import * as fs from 'fs'
import { parse } from 'jsonc-parser'
import { TAuthContext } from '@bearer/functions/lib/declaration'

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

async function getFunctionContext(locator: LocationProvider) {
  const localConfig = locator.localConfigPath
  let context = {} as TAuthContext
  if (fs.existsSync(localConfig)) {
    const rawConfig = fs.readFileSync(localConfig, { encoding: 'utf8' })
    const parsed = parse(rawConfig)
    const { setup } = parsed || { setup: null }
    debug('local config: %j', parsed)
    if (setup && setup.auth) {
      context.authAccess = setup.auth
    }
  } else {
    debug('no local config found')
  }

  return context
}

export default function startLocalDevelopmentServer(
  emitter,
  config: Config,
  locator: LocationProvider,
  { logs = true, force = false } = {}
) {
  const rootLevel = locator.integrationRoot

  const router = new Router({ prefix: '/api/' })

  return new Promise<string>(async (resolve, reject) => {
    try {
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
      const expectedPort = process.env.PORT ? parseInt(process.env.PORT) : DEFAULT_PORT

      const port = await getPort({ port: expectedPort })

      // tslint:disable-next-line:no-http-string
      const host = `http://localhost:${port}`

      if (expectedPort !== DEFAULT_PORT) {
        emitter.emit('start:localServer:customPort', { port, host })
      }

      if (port !== Number(expectedPort) && !force) {
        throw `Could not start local server port ${expectedPort} is already in use. You can specify your own port by running PORT=3322 yarn bearer start`
      }
      const bearerBaseURL = `${host}/`
      process.env.bearerBaseURL = bearerBaseURL

      const context = await getFunctionContext(locator)
      router.post(
        `v4/backend/functions/${config.buid}/:functionName`,
        functionHandler(distPath, context, bearerBaseURL, true),
        (ctx, _next) => ctx.ok(ctx.funcDatum)
      )

      router.post(
        `v4/functions/${config.buid}/:functionName`,
        functionHandler(distPath, context, bearerBaseURL),
        (ctx, _next) => ctx.ok(ctx.funcDatum)
      )

      router.post(`backend/:functionName`, functionHandler(distPath, context, bearerBaseURL, true), (ctx, _next) =>
        ctx.ok(ctx.funcDatum)
      )

      router.post(`:functionName`, functionHandler(distPath, context, bearerBaseURL), (ctx, _next) =>
        ctx.ok(ctx.funcDatum)
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

const functionHandler = (distPath: string, devFunctionsContext, bearerBaseURL: string, isBackend = false) => async (
  ctx,
  next
) =>
  new Promise(async (resolve, _reject) => {
    try {
      const func = requireUncached(`${distPath}/${ctx.params.functionName}`).default

      const userDefinedData = await loadUserDefinedData({ query: ctx.query })

      const datum = await func.init()(
        {
          context: {
            ...devFunctionsContext,
            bearerBaseURL,
            ...userDefinedData,
            isBackend
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
