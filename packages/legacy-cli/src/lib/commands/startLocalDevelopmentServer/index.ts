import * as chokidar from 'chokidar'
import * as cosmiconfig from 'cosmiconfig'
import * as getPort from 'get-port'
import * as Logger from 'koa-logger'
import * as Router from 'koa-router'

import { transpileIntents } from '../../buildArtifact'
import LocationProvider from '../../locationProvider'
import { Config } from '../../types'

import auth from './auth'
import server = require('./server')
import Storage from './storage'
import { loadUserDefinedData } from './utils'

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
  const rootLevel = locator.scenarioRoot

  const LOCAL_DEV_CONFIGURATION = 'dev'
  const explorer = cosmiconfig(LOCAL_DEV_CONFIGURATION, {
    searchPlaces: [`config.${LOCAL_DEV_CONFIGURATION}.js`]
  })
  const router = new Router({ prefix: '/api/' })

  return new Promise(async (resolve, reject) => {
    try {
      const { config: devIntentsContext = {} } = (await explorer.search(rootLevel)) || {}
      const distPath = locator.buildIntentsResourcePath('dist')

      // tslint:disable-next-line:no-inner-declarations
      async function refreshIntents() {
        try {
          emitter.emit('start:localServer:generatingIntents:start')
          await transpileIntents(locator.srcIntentsDir, distPath)
          emitter.emit('start:localServer:generatingIntents:stop')
        } catch (error) {
          emitter.emit('start:localServer:generatingIntents:failed', { error })
        }
      }
      await refreshIntents()

      chokidar
        .watch('.', {
          ignored: /(^|[\/\\])\../,
          cwd: locator.srcIntentsDir,
          ignoreInitial: true,
          persistent: true,
          followSymlinks: false
        })
        .on('add', refreshIntents)
        .on('change', refreshIntents)

      const storage = Storage()
      const port = await getPort({ port: 3000 })
      // tslint:disable-next-line:no-http-string
      const bearerBaseURL = `http://localhost:${port}/`
      process.env.bearerBaseURL = bearerBaseURL

      router.post(
        `v2/intents/${config.scenarioUuid}/:intentName`,
        intentHandler(distPath, devIntentsContext, bearerBaseURL),
        (ctx, _next) => {
          if (ctx.intentDatum.error) {
            ctx.badRequest({ error: ctx.intentDatum.error })
          } else {
            ctx.ok(ctx.intentDatum)
          }
        }
      )

      router.all(
        `v1/${config.scenarioUuid}/:intentName`,
        intentHandler(distPath, devIntentsContext, bearerBaseURL),
        (ctx, _next) => {
          if (ctx.intentDatum.error) {
            ctx.badRequest({ error: ctx.intentDatum.error })
          } else {
            ctx.ok(ctx.intentDatum)
          }
        }
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

const intentHandler = (distPath: string, devIntentsContext, bearerBaseURL: string) => async (ctx, next) =>
  new Promise(async (resolve, _reject) => {
    try {
      const intent = requireUncached(`${distPath}/${ctx.params.intentName}`).default

      const userDefinedData = await loadUserDefinedData({ query: ctx.query })

      intent.intentType.intent(intent.action)(
        {
          context: {
            ...devIntentsContext.global,
            ...devIntentsContext[ctx.params.intentName],
            bearerBaseURL,
            ...userDefinedData
          },
          queryStringParameters: ctx.query,
          body: JSON.stringify(ctx.request.body)
        },
        {},
        (_err, datum) => {
          ctx.intentDatum = datum
          next()
          resolve()
        }
      )
    } catch (e) {
      console.log('ERROR: ', e)
      if (e.code === 'MODULE_NOT_FOUND') {
        ctx.intentDatum = { error: `Intent '${ctx.params.intentName}' Not Found` }
      } else {
        ctx.intentDatum = { error: e }
      }
      await next()
      resolve()
    }
  })
