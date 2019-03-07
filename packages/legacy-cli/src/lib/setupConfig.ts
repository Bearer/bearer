const findUp = require('find-up')
const rc = require('rc')
const os = require('os')

const fs = require('fs')
const ini = require('ini')
const path = require('path')
const { spawnSync } = require('child_process')

import { BaseConfig, BearerConfig, BearerEnv, Config, IntegrationConfig } from './types'

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
    get bearerConfig(): BearerConfig {
      return rc('bearer')
    },
    get integrationConfig(): IntegrationConfig {
      return rc('integration')
    },
    get orgId(): string | undefined {
      return this.integrationConfig.orgId
    },
    get integrationTitle(): string | undefined {
      return this.integrationConfig.integrationTitle
    },
    get integrationId(): string | undefined {
      return this.integrationConfig.integrationId
    },
    get integrationUuid(): string {
      if (this.hasIntegrationLinked) {
        return `${this.orgId}-${this.integrationId}`
      }
      return 'unset-scenari0-uuid'
    },
    get hasIntegrationLinked(): boolean {
      return Boolean(this.orgId) && Boolean(this.integrationId)
    },
    get rootPathRc(): string {
      return findUp.sync('.integrationrc')
    },
    setIntegrationConfig(config: { integrationTitle: string; orgId: string; integrationId: string }) {
      const { integrationTitle, orgId, integrationId } = config
      fs.writeFileSync(this.rootPathRc, ini.stringify({ integrationTitle, orgId, integrationId }))
    },
    storeBearerConfig(config) {
      const { Username, ExpiresAt, authorization } = config
      fs.writeFileSync(
        this.bearerConfig.config || path.join(os.homedir(), '.bearerrc'),
        ini.stringify({
          Username,
          ExpiresAt,
          authorization
        })
      )
    }
  }
}
