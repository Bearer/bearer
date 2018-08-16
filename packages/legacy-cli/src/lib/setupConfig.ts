const rc = require('rc')
const os = require('os')
const findUp = require('find-up')
const fs = require('fs')
const ini = require('ini')
const path = require('path')
const { spawnSync } = require('child_process')

import { BaseConfig, BearerConfig, BearerEnv, Config, ScenarioConfig } from './types'

const configs: Record<BearerEnv, BaseConfig> = {
  dev: {
    DeploymentUrl: 'https://developer.dev.bearer.sh/v1/',
    IntegrationServiceHost: 'https://int.dev.bearer.sh/',
    IntegrationServiceUrl: 'https://int.dev.bearer.sh/api/v1/',
    DeveloperPortalAPIUrl: 'https://app.staging.bearer.sh/graphql',
    DeveloperPortalUrl: 'https://app.staging.bearer.sh/',
    CdnHost: 'https://static.dev.bearer.sh',
    BearerEnv: 'dev'
  },
  staging: {
    DeploymentUrl: 'https://developer.staging.bearer.sh/v1/',
    IntegrationServiceHost: 'https://int.staging.bearer.sh/',
    IntegrationServiceUrl: 'https://int.staging.bearer.sh/api/v1/',
    DeveloperPortalAPIUrl: 'https://app.staging.bearer.sh/graphql',
    DeveloperPortalUrl: 'https://app.staging.bearer.sh/',
    CdnHost: 'https://static.staging.bearer.sh',
    BearerEnv: 'staging'
  },
  production: {
    DeploymentUrl: 'https://developer.bearer.sh/v1/',
    IntegrationServiceHost: 'https://int.bearer.sh/',
    IntegrationServiceUrl: 'https://int.bearer.sh/api/v1/',
    DeveloperPortalAPIUrl: 'https://app.bearer.sh/graphql',
    DeveloperPortalUrl: 'https://app.bearer.sh/',
    CdnHost: 'https://static.bearer.sh',
    BearerEnv: 'production'
  }
}

export default (): Config => {
  const { BEARER_ENV = 'production' } = process.env
  const setup: BaseConfig = configs[BEARER_ENV]

  const isYarnInstalled = !!spawnSync('yarn', ['bin']).output

  return {
    ...setup,
    isYarnInstalled,
    command: isYarnInstalled ? 'yarn' : 'npm',
    get bearerConfigFileName(): string {
      if (setup.BearerEnv === 'production') {
        return 'bearer'
      } else {
        return `${setup.BearerEnv}.bearer`
      }
    },
    get bearerConfig(): BearerConfig {
      return rc(this.bearerConfigFileName)
    },
    get scenarioConfig(): ScenarioConfig {
      return rc(this.rootPathFileName)
    },
    get orgId(): string | undefined {
      return this.scenarioConfig.orgId
    },
    get scenarioTitle(): string | undefined {
      return this.scenarioConfig.scenarioTitle
    },
    get scenarioId(): string | undefined {
      return this.scenarioConfig.scenarioId
    },
    get scenarioUuid(): string | undefined {
      if (!this.orgId || !this.scenarioId) {
        return undefined
      }
      return `${this.orgId}-${this.scenarioId}`
    },
    get rootPathFileName(): string {
      if (setup.BearerEnv === 'production') {
        return 'scenario'
      } else {
        return `${setup.BearerEnv}.scenario`
      }
    },
    get rootPathRc(): string {
      return findUp.sync(`.${this.rootPathFileName}rc`)
    },
    setScenarioConfig(config: { scenarioTitle: string; orgId: string; scenarioId: string }) {
      const { scenarioTitle, orgId, scenarioId } = config
      fs.writeFileSync(this.rootPathRc, ini.stringify({ scenarioTitle, orgId, scenarioId }))
    },
    storeBearerConfig(config) {
      const { Username, ExpiresAt, authorization } = config
      fs.writeFileSync(
        this.bearerConfig.config || path.join(os.homedir(), `${this.bearerConfigFileName}.rc`),
        ini.stringify({
          Username,
          ExpiresAt,
          authorization
        })
      )
    }
  }
}
