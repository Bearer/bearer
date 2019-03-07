import { spawnSync } from 'child_process'
import * as findUp from 'find-up'
import * as fs from 'fs'
import * as ini from 'ini'
import * as os from 'os'
import * as path from 'path'
import * as rc from 'rc'
import { promisify } from 'util'

import { BaseConfig, BearerConfig, BearerEnv, Config, IntegrationConfig } from '../types'

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

export default (runPath: string = process.cwd()): Config => {
  const { BEARER_ENV = 'production' } = process.env
  const setup: BaseConfig = configs[BEARER_ENV as BearerEnv]

  const isYarnInstalled = !!spawnSync('yarn', ['bin']).output
  const integrationLocation = runPath.startsWith('~')
    ? path.resolve(runPath.replace(/^~/, os.homedir()))
    : path.resolve(runPath)
  return {
    ...setup,
    isYarnInstalled,
    runPath: integrationLocation,
    command: isYarnInstalled ? 'yarn' : 'npm',
    get isIntegrationLocation(): boolean {
      return this.rootPathRc !== null
    },
    get bearerConfig(): BearerConfig {
      return rc('bearer')
    },
    get integrationConfig(): IntegrationConfig {
      return rc('integration', { config: path.join(integrationLocation, '.integrationrc') })
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
      return 'unset-integration-uuid'
    },
    get hasIntegrationLinked(): boolean {
      return Boolean(this.orgId) && Boolean(this.integrationId)
    },
    get rootPathRc(): string | null {
      return findUp.sync('.integrationrc', { cwd: integrationLocation })
    },
    setIntegrationConfig(config: { integrationTitle: string; orgId: string; integrationId: string }) {
      const { integrationTitle, orgId, integrationId } = config
      if (this.rootPathRc) {
        fs.writeFileSync(this.rootPathRc, ini.stringify({ integrationTitle, orgId, integrationId }))
      }
    },
    async storeBearerConfig(config) {
      const { Username, ExpiresAt, authorization } = config
      const writeFile = promisify(fs.writeFile)
      try {
        await writeFile(
          this.bearerConfig.config || path.join(os.homedir(), '.bearerrc'),
          ini.stringify({
            Username,
            ExpiresAt,
            authorization
          })
        )
      } catch (e) {
        console.error('Error while writing the token', e)
      }
      this.bearerConfig.authorization = authorization
    }
  }
}
