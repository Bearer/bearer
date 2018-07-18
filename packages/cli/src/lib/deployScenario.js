const { spawn } = require('child_process')
const pathJs = require('path')
const fs = require('fs')
const util = require('util')
const exec = util.promisify(require('child_process').exec)
const { prepare } = require('./commands/startCommand')
const buildArtifact = require('./buildArtifact')
const pushScenario = require('./pushScenario')
const pushScreens = require('./pushScreens')
const assembly = require('./assemblyScenario')
const refreshToken = require('./refreshToken')
const invalidateCloudFront = require('./invalidateCloudFront')
const developerPortal = require('./developerPortal')

function buildIntents(rootLevel, scenarioUuid, emitter, config) {
  return new Promise(async (resolve, reject) => {
    const artifactDirectory = pathJs.join(rootLevel, '.bearer')
    const intentsDirectory = pathJs.join(rootLevel, 'intents')

    if (!fs.existsSync(artifactDirectory)) {
      fs.mkdirSync(artifactDirectory)
    }
    try {
      const scenarioArtifact = pathJs.join(
        artifactDirectory,
        `${scenarioUuid}.zip`
      )
      const handler = pathJs.join(artifactDirectory, config.HandlerBase)
      const output = fs.createWriteStream(scenarioArtifact)

      emitter.emit('intents:installingDependencies')
      await exec('yarn install', { cwd: intentsDirectory })

      await buildArtifact(
        output,
        handler,
        { path: intentsDirectory, scenarioUuid },
        emitter
      )
      return resolve(scenarioArtifact)
    } catch (e) {
      return reject(e)
    }
  })
}

function transpileStep(
  emitter,
  screensDirectory,
  scenarioUuid,
  integrationHost
) {
  return new Promise(async (resolve, reject) => {
    emitter.emit('start:prepare:transpileStep')
    const bearerTranspiler = spawn(
      'node',
      [pathJs.join(__dirname, 'startTranspiler.js'), '--no-watcher'],
      {
        cwd: screensDirectory,
        env: {
          ...process.env,
          BEARER_SCENARIO_ID: scenarioUuid,
          BEARER_INTEGRATION_HOST: integrationHost
        },
        stdio: ['pipe', 'pipe', 'pipe', 'ipc']
      }
    )

    bearerTranspiler.on('close', (...args) => {
      emitter.emit('start:prepare:transpileStep:close', args)
      resolve(...args)
    })
    bearerTranspiler.stderr.on('data', (...args) => {
      emitter.emit('start:prepare:transpileStep:command:error', args)
      reject(...args)
    })
    bearerTranspiler.on('message', ({ event }) => {
      if (event === 'transpiler:initialized') {
        emitter.emit('start:prepare:transpileStep:done')
        resolve(bearerTranspiler)
      } else {
        emitter.emit('start:prepare:transpileStep:error')
        reject(event)
      }
    })
  })
}

const deployIntents = ({ scenarioUuid }, emitter, config) =>
  new Promise(async (resolve, reject) => {
    const { rootPathRc } = config

    if (!rootPathRc) {
      emitter.emit('rootPath:doesntExist')
      process.exit(1)
    }

    const rootLevel = pathJs.dirname(rootPathRc)

    try {
      const scenarioArtifact = await buildIntents(
        rootLevel,
        scenarioUuid,
        emitter,
        config
      )
      await pushScenario(
        scenarioArtifact,
        {
          Key: scenarioUuid
        },
        emitter,
        config
      )

      await assembly(scenarioUuid, emitter, config)
      resolve()
    } catch (e) {
      reject(e)
    }
  })

const deployScreens = ({ scenarioUuid }, emitter, config) =>
  new Promise(async (resolve, reject) => {
    const {
      scenarioConfig: { scenarioTitle },
      bearerConfig: { OrgId },
      BearerEnv
    } = config

    try {
      const { buildDirectory } = await prepare(emitter, config)()
      if (!buildDirectory) {
        process.exit(1)
        return false
      }

      await transpileStep(
        emitter,
        pathJs.join(buildDirectory, '..'),
        scenarioUuid,
        config.IntegrationServiceHost
      )

      emitter.emit('screens:generateSetupComponent')

      await exec('bearer generate --setup', { cwd: buildDirectory })

      emitter.emit('screens:generateConfigComponent')
      await exec('bearer generate --config', { cwd: buildDirectory })

      emitter.emit('screens:buildingDist')
      await exec('yarn build', {
        cwd: buildDirectory,
        pwd: buildDirectory,
        env: {
          BEARER_SCENARIO_ID: scenarioUuid,
          ...process.env,
          CDN_HOST: `https://static.${BearerEnv}.bearer.sh/${OrgId}/${scenarioTitle}/dist/${scenarioTitle}/`
        }
      })

      emitter.emit('screens:pushingDist')
      await pushScreens(buildDirectory, scenarioTitle, OrgId, emitter, config)

      emitter.emit('screen:upload:success')
      await invalidateCloudFront(emitter, config)
      return resolve()
    } catch (e) {
      emitter.emit('deployScenario:deployScreens:error', e)
      return reject(e)
    }
  })

module.exports = {
  deployIntents,
  deployScreens,
  buildIntents,
  deployScenario: ({ scenarioUuid }, emitter, config) =>
    new Promise(async (resolve, reject) => {
      let calculatedConfig = config

      try {
        const { ExpiresAt } = config.bearerConfig

        if (ExpiresAt < Date.now()) {
          calculatedConfig = await refreshToken(config, emitter)
        }
        await developerPortal(emitter, 'predeploy', calculatedConfig)
        await deployIntents({ scenarioUuid }, emitter, calculatedConfig)
        await deployScreens({ scenarioUuid }, emitter, calculatedConfig)
        await developerPortal(emitter, 'deployed', calculatedConfig)
        resolve()
      } catch (e) {
        await developerPortal(emitter, 'cancelled', calculatedConfig)
        reject(e)
      }
    })
}
