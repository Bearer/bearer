import * as getPort from 'get-port'
import * as Router from 'koa-router'
import * as cosmiconfig from 'cosmiconfig'
import * as Logger from 'koa-logger'
import * as chokidar from 'chokidar'

import { transpileIntents } from '../../buildArtifact'
import server = require('./server')
import Storage from './storage'
import auth from './auth'
import LocationProvider from '../../locationProvider'
import { Config } from '../../types'

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
  const router = new Router({ prefix: '/api/v1/' })

  return new Promise(async (resolve, reject) => {
    try {
      const { config: devIntentsContext = {} } = (await explorer.search(rootLevel)) || {}
      const distPath = locator.buildIntentsResourcePath('dist')

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

      const port = await getPort({ port: 3000 })
      const bearerBaseURL = `http://localhost:${port}/`
      router.all(
        `${config.scenarioUuid}/:intentName`,
        (ctx, next) =>
          new Promise((resolve, reject) => {
            try {
              const intent = requireUncached(`${distPath}/${ctx.params.intentName}`).default
              intent.intentType.intent(intent.action)(
                {
                  context: {
                    ...devIntentsContext.global,
                    ...devIntentsContext[ctx.params.intentName],
                    bearerBaseURL
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
              reject({ error: e.toString() })
            }
          }),
        ctx => ctx.ok(ctx.intentDatum)
      )

      const storage = Storage()
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
