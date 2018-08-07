import * as getPort from 'get-port'
import * as Router from 'koa-router'
import * as unzip from 'unzip-stream'
import * as fs from 'fs-extra'
import * as cosmiconfig from 'cosmiconfig'
import * as Logger from 'koa-logger'

import server = require('./server')
import Storage from './storage'
import auth from './auth'
import { buildIntents } from '../../deployScenario'
import LocationProvider from '../../locationProvider'

function startLocalDevelopmentServer(_scenarioUuid, emitter, config, locator: LocationProvider, logs: boolean = true) {
  const rootLevel = locator.scenarioRoot
  const buildDir = locator.buildIntentsDir

  const LOCAL_DEV_CONFIGURATION = 'dev'
  const explorer = cosmiconfig(LOCAL_DEV_CONFIGURATION, {
    searchPlaces: [`config.${LOCAL_DEV_CONFIGURATION}.js`]
  })
  const router = new Router({ prefix: '/api/v1/' })

  return new Promise(async (resolve, reject) => {
    try {
      const { config: devIntentsContext = {} } = (await explorer.search(rootLevel)) || {}
      const intentsArtifact = await buildIntents(emitter, config, locator)

      fs.ensureDirSync(buildDir)
      await new Promise((resolve, reject) => {
        fs.createReadStream(intentsArtifact)
          .pipe(unzip.Extract({ path: buildDir }))
          .on('close', resolve)
          .on('error', reject)
      })

      const { integration_uuid } = require(locator.buildIntentsResourcePath('bearer.config.json'))

      const port = await getPort({ port: 3000 })
      const bearerBaseURL = `http://localhost:${port}/`
      router.all(
        `${integration_uuid}/:intentName`,
        (ctx, next) =>
          new Promise((resolve, _reject) => {
            const intent = require(`${buildDir}/${ctx.params.intentName}`).default
            intent.intentType.intent(intent.action)(
              {
                context: {
                  ...devIntentsContext.global,
                  ...devIntentsContext[ctx.params.intentName],
                  bearerBaseURL
                },
                queryStringParameters: ctx.query,
                body: ctx.request.body
              },
              {},
              (_err, datum) => {
                ctx.intentDatum = datum
                next()
                resolve()
              }
            )
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

module.exports = startLocalDevelopmentServer
