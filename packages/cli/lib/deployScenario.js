const pathJs = require('path')
const fs = require('fs')
const util = require('util')
const exec = util.promisify(require('child_process').exec)

const { prepare } = require('./commands/startCommand')
const buildArtifact = require('./buildArtifact')
const pushScenario = require('./pushScenario')
const pushScreens = require('./pushScreens')
const assembly = require('./assemblyScenario')
const storeCredentials = require('./storeCredentials')
const refreshToken = require('./refreshToken')
const invalidateCloudFront = require('./invalidateCloudFront')

const AUTH_CONFIG_FILE = 'auth.config.json'

module.exports = async ({ scenarioUuid }, emitter, config) => {
  const {
    rootPathRc,
    scenarioConfig: { scenarioTitle },
    bearerConfig: { OrgId },
    BearerEnv
  } = config

  if (!rootPathRc) {
    emitter.emit('rootPath:doesntExist')
    process.exit(1)
  }

  const rootLevel = pathJs.dirname(rootPathRc)
  const artifactDirectory = pathJs.join(rootLevel, '.bearer')
  const intentsDirectory = pathJs.join(rootLevel, 'intents')

  if (!fs.existsSync(artifactDirectory)) {
    fs.mkdirSync(artifactDirectory)
  }

  const scenarioArtifact = pathJs.join(artifactDirectory, `${scenarioUuid}.zip`)
  const handler = pathJs.join(artifactDirectory, config.HandlerBase)
  const output = fs.createWriteStream(scenarioArtifact)

  try {
    emitter.emit('intents:installingDependencies')
    await exec('yarn install', { cwd: intentsDirectory })

    await buildArtifact(
      output,
      handler,
      { path: intentsDirectory, scenarioUuid },
      emitter
    )

    const authConfigFilePath = pathJs.join(intentsDirectory, AUTH_CONFIG_FILE)
    await storeCredentials(authConfigFilePath, config, emitter)

    const { ExpiresAt } = config.bearerConfig
    let calculatedConfig = config

    if (ExpiresAt < Date.now()) {
      calculatedConfig = await refreshToken(config, emitter)
    }

    await pushScenario(
      scenarioArtifact,
      {
        Key: scenarioUuid
      },
      emitter,
      calculatedConfig
    )

    await assembly(scenarioUuid, emitter, calculatedConfig)

    const { buildDirectory: screensDirectory } = await prepare(
      emitter,
      config
    )()
    if (!screensDirectory) {
      process.exit(1)
      return false
    }

    emitter.emit('screens:generateSetupComponent')
    await exec('bearer generate --setup', { cwd: screensDirectory })

    emitter.emit('screens:generateConfigComponent')
    await exec('bearer generate --config', { cwd: screensDirectory })

    emitter.emit('screens:buildingDist')
    await exec('yarn build', {
      cwd: screensDirectory,
      pwd: screensDirectory,
      env: {
        BEARER_SCENARIO_ID: scenarioUuid,
        ...process.env,
        CDN_HOST: `https://static.${BearerEnv}.bearer.sh/${OrgId}/${scenarioTitle}/dist/${scenarioTitle}/`
      }
    })

    emitter.emit('screens:pushingDist')
    await pushScreens(
      screensDirectory,
      scenarioTitle,
      OrgId,
      emitter,
      calculatedConfig
    )

    emitter.emit('screen:upload:success')
    await invalidateCloudFront(emitter, calculatedConfig)
  } catch (e) {
    console.log(e)
  }
}
