import { exec, spawn } from 'child_process'
import * as fs from 'fs-extra'
import * as path from 'path'
import { promisify } from 'util'

import * as assembly from './assemblyScenario'
import buildArtifact from './buildArtifact'
import { prepare } from './commands/startCommand'
import * as developerPortal from './developerPortal'
import invalidateCloudFront from './invalidateCloudFront'
import LocationProvider from './locationProvider'
import * as pushScenario from './pushScenario'
import * as pushViews from './pushViews'

import { Config } from './types'

const execPromise = promisify(exec)

export function buildIntents(emitter, config: Config, locator: LocationProvider) {
  return new Promise(async (resolve, reject) => {
    const { scenarioUuid } = config
    const intentsDirectory = locator.srcIntentsDir

    try {
      emitter.emit('intents:installingDependencies')
      // TODOs: use root node modules
      await execPromise(`${config.command} install`, { cwd: intentsDirectory })

      const scenarioArtifact = locator.buildArtifactResourcePath(`${scenarioUuid}.zip`)
      const output = fs.createWriteStream(scenarioArtifact)
      buildArtifact(output, { scenarioUuid }, emitter, locator)
        .then(() => {
          emitter.emit('intents:buildIntents:succeeded')
          resolve(scenarioArtifact)
        })
        .catch(error => {
          emitter.emit('intents:buildIntents:failed', { error })
          reject(error)
        })
    } catch (e) {
      reject(e)
    }
  })
}

export function deployIntents(emitter, config: Config, locator: LocationProvider) {
  return new Promise(async (resolve, _reject) => {
    const { rootPathRc } = config

    if (!rootPathRc) {
      emitter.emit('rootPath:doesntExist')
      process.exit(1)
    }

    try {
      const scenarioArtifact = await buildIntents(emitter, config, locator)
      await pushScenario(scenarioArtifact, emitter, config)

      await assembly(emitter, config)
      resolve()
    } catch (e) {
      emitter.emit('deployIntents:error', e)
      resolve()
    }
  })
}

export function deployViews(emitter, config: Config, locator: LocationProvider) {
  return new Promise(async (resolve, reject) => {
    const { orgId, scenarioUuid, scenarioId, CdnHost } = config

    try {
      const { buildDirectory } = await prepare(emitter, config, locator)({
        install: true,
        watchMode: null
      })
      if (!buildDirectory) {
        process.exit(1)
        return false
      }

      await transpileStep(emitter, locator, config)

      emitter.emit('views:buildingDist')
      await execPromise(`${config.command} build`, {
        cwd: buildDirectory,
        env: {
          BEARER_SCENARIO_ID: scenarioUuid,
          BEARER_SCENARIO_TAG_NAME: scenarioId,
          BEARER_INTEGRATION_HOST: config.IntegrationServiceHost,
          BEARER_AUTHORIZATION_HOST: config.IntegrationServiceHost,
          ...process.env,
          CDN_HOST: `${CdnHost}/${orgId}/${scenarioId}/dist/${scenarioId}/`
        }
      })

      emitter.emit('views:pushingDist')
      await pushViews(buildDirectory, emitter, config)

      emitter.emit('view:upload:success')
      await invalidateCloudFront(emitter, config)
      resolve()
    } catch (e) {
      emitter.emit('deployScenario:deployViews:error', e)
      console.error(e)
      reject(e)
    }
  })
}

function transpileStep(emitter, locator: LocationProvider, config: Config) {
  return new Promise(async (resolve, reject) => {
    const { scenarioUuid, IntegrationServiceHost, scenarioId, orgId } = config
    emitter.emit('start:prepare:transpileStep')
    const prefix = ['bearer', scenarioId].join('-')
    const suffix = orgId
    const options = [
      path.join(__dirname, 'startTranspiler.js'),
      '--no-watcher',
      '--prefix-tag',
      prefix,
      '--suffix-tag',
      suffix
    ]
    const bearerTranspiler = spawn('node', options, {
      cwd: locator.scenarioRoot,
      env: {
        ...process.env,
        BEARER_SCENARIO_TAG_NAME: scenarioId,
        BEARER_SCENARIO_ID: scenarioUuid,
        BEARER_INTEGRATION_HOST: IntegrationServiceHost,
        BEARER_AUTHORIZATION_HOST: IntegrationServiceHost
      },
      stdio: ['pipe', 'pipe', 'pipe', 'ipc']
    })

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

export interface IDeployOptions {
  noViews?: boolean
  noIntents?: boolean
}

export function deployScenario(
  { noViews = false, noIntents = false }: IDeployOptions,
  emitter,
  config: Config,
  locator
) {
  return new Promise(async (resolve, reject) => {
    try {
      await developerPortal(emitter, 'predeploy', config)
      if (!noIntents) {
        await deployIntents(emitter, config, locator)
      }
      if (!noViews) {
        await deployViews(emitter, config, locator)
      }
      await developerPortal(emitter, 'deployed', config)
      resolve()
    } catch (e) {
      emitter.emit('deployScenario:error', e)
      await developerPortal(emitter, 'cancelled', config)
      reject(e)
    }
  })
}
