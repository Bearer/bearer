import { spawnSync } from 'child_process'
import * as findUp from 'find-up'
import * as fs from 'fs'
import * as ini from 'ini'
import * as os from 'os'
import * as path from 'path'
import * as rc from 'rc'

import { BaseConfig, BearerConfig, BearerEnv, Config, ScenarioConfig } from '../types'

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
  const setup: BaseConfig = configs[BEARER_ENV as BearerEnv]

  const isYarnInstalled = !!spawnSync('yarn', ['bin']).output

  return {
    ...setup,
    isYarnInstalled,
    command: isYarnInstalled ? 'yarn' : 'npm',
    get isScenarioLocation(): boolean {
      return this.rootPathRc !== null
    },
    get bearerConfig(): BearerConfig {
      return rc('bearer')
    },
    get scenarioConfig(): ScenarioConfig {
      return rc('scenario')
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
    get scenarioUuid(): string {
      if (this.hasScenarioLinked) {
        return `${this.orgId}-${this.scenarioId}`
      }
      return 'unset-scenario-uuid'
    },
    get hasScenarioLinked(): boolean {
      return Boolean(this.orgId) && Boolean(this.scenarioId)
    },
    get rootPathRc(): string | null {
      return findUp.sync('.scenariorc')
    },
    setScenarioConfig(config: { scenarioTitle: string; orgId: string; scenarioId: string }) {
      const { scenarioTitle, orgId, scenarioId } = config
      if (this.rootPathRc) {
        fs.writeFileSync(this.rootPathRc, ini.stringify({ scenarioTitle, orgId, scenarioId }))
      }
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
