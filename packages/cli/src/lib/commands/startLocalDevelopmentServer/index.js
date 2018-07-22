const path = require('path')
const server = require('./server')
const getPort = require('get-port')
const Router = require('koa-router')
const unzip = require('unzip')
const fs = require('fs-extra')
const cosmiconfig = require('cosmiconfig')
const storage = require('./storage')

const LOCAL_DEV_CONFIGURATION = 'dev'
const explorer = cosmiconfig(LOCAL_DEV_CONFIGURATION)

const router = new Router({ prefix: '/api/v1/' })

async function startLocalDevelopmentServer(
  rootLevel,
  scenarioUuid,
  emitter,
  config
) {
  return new Promise(async (resolve, reject) => {
    try {
      const { config: devIntentsContext } =
        (await explorer.search(rootLevel)) || {}
      const { buildIntents } = require(path.join(
        __dirname,
        '..',
        '..',
        'deployScenario'
      ))
      const intentsArtifact = await buildIntents(
        rootLevel,
        scenarioUuid,
        emitter,
        config
      )

      const buildDir = path.join(rootLevel, 'intents', '.build')
      fs.ensureDirSync(buildDir)

      await new Promise((resolve, reject) => {
        fs.createReadStream(intentsArtifact)
          .pipe(unzip.Extract({ path: buildDir }))
          .on('close', resolve)
          .on('error', reject)
      })
      const lambdas = require(buildDir)

      const { integration_uuid, intents } = require(path.join(
        buildDir,
        'bearer.config.json'
      ))

      const port = await getPort({ port: 3000 })
      const bearerBaseURL = `http://localhost:${port}/`
      for (let intent of intents) {
        const intentName = Object.keys(intent)[0]
        const endpoint = `${integration_uuid}/${intentName}`
        router.all(
          endpoint,
          (ctx, next) =>
            new Promise((resolve, reject) => {
              console.log('Body is: ', ctx.request.body)
              lambdas[intentName](
                {
                  context: {
                    ...devIntentsContext.global,
                    ...devIntentsContext[intentName],
                    bearerBaseURL
                  },
                  queryStringParameters: ctx.query,
                  body: ctx.request.body
                },
                {},
                (err, datum) => {
                  ctx.intentDatum = datum
                  next()
                  resolve()
                }
              )
            }),
          ctx => ctx.ok(ctx.intentDatum)
        )
      }

      server.use(router.routes())
      server.use(router.allowedMethods())
      server.use(storage.routes())
      server.use(storage.allowedMethods())

      server.listen(port, () => {
        emitter.emit('start:localServer:start', { port })
        emitter.emit('start:localServer:endpoints', {
          endpoints: router.stack
        })
      })

      resolve(bearerBaseURL)
    } catch (e) {
      reject(e)
    }
  })
}

module.exports = startLocalDevelopmentServer
