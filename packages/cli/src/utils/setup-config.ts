import { spawnSync } from 'child_process'
import * as findUp from 'find-up'
import * as fs from 'fs'
import * as ini from 'ini'
import * as os from 'os'
import * as path from 'path'
import * as rc from 'rc'
import { promisify } from 'util'

import { BaseConfig, BearerConfig, BearerEnv, IntegrationConfig, TAccessToken } from '../types'

const writeFile = promisify(fs.writeFile)
const readFile = promisify(fs.readFile)

const configs: Record<BearerEnv, BaseConfig> = {
  dev: {
    DeploymentUrl: 'https://developer.dev.bearer.sh/v1/',
    IntegrationServiceHost: 'https://int.dev.bearer.sh/',
    IntegrationServiceUrl: 'https://int.dev.bearer.sh/api/v1/',
    DeveloperPortalAPIUrl: 'https://api.staging.bearer.sh/graphql',
    DeveloperPortalUrl: 'https://app.staging.bearer.sh/',
    CdnHost: 'https://static.dev.bearer.sh',
    BearerEnv: 'dev',
    LoginDomain: 'https://login.bearer.sh'
  },
  staging: {
    DeploymentUrl: 'https://developer.staging.bearer.sh/cicd/v2',
    IntegrationServiceHost: 'https://int.staging.bearer.sh/',
    IntegrationServiceUrl: 'https://int.staging.bearer.sh/api/v1/',
    DeveloperPortalAPIUrl: 'https://api.staging.bearer.sh/graphql',
    DeveloperPortalUrl: 'https://app.staging.bearer.sh/',
    CdnHost: 'https://static.staging.bearer.sh',
    BearerEnv: 'staging',
    LoginDomain: 'https://login.bearer.sh'
  },
  production: {
    DeploymentUrl: 'https://developer.bearer.sh/cicd/v2/',
    IntegrationServiceHost: 'https://int.bearer.sh/',
    IntegrationServiceUrl: 'https://int.bearer.sh/api/v1/',
    DeveloperPortalAPIUrl: 'https://api.bearer.sh/graphql',
    DeveloperPortalUrl: 'https://app.bearer.sh/',
    CdnHost: 'https://static.bearer.sh',
    BearerEnv: 'production',
    LoginDomain: 'https://login.bearer.sh'
  }
}

export class Config {
  integrationLocation: string

  constructor(readonly runPath: string) {
    this.integrationLocation = this.runPath.startsWith('~')
      ? path.resolve(runPath.replace(/^~/, os.homedir()))
      : path.resolve(runPath)
  }

  get command(): 'yarn' | 'npm' {
    return this.isYarnInstalled ? 'yarn' : 'npm'
  }

  get isIntegrationLocation(): boolean {
    return this.rootPathRc !== null
  }

  get bearerConfig(): BearerConfig {
    return rc('bearer')
  }

  get integrationConfig(): IntegrationConfig {
    return rc('integration', { config: path.join(this.integrationLocation, '.integrationrc') })
  }

  get orgId(): string | undefined {
    return this.integrationConfig.orgId
  }

  get integrationTitle(): string | undefined {
    return this.integrationConfig.integrationTitle
  }

  get integrationId(): string | undefined {
    return this.integrationConfig.integrationId
  }

  get integrationUuid(): string {
    if (this.hasIntegrationLinked) {
      return `${this.orgId}-${this.integrationId}`
    }
    return 'unset-integration-uuid'
  }

  get hasIntegrationLinked(): boolean {
    return Boolean(this.orgId) && Boolean(this.integrationId)
  }

  get rootPathRc(): string | null {
    return findUp.sync('.integrationrc', { cwd: this.integrationLocation })
  }

  get isYarnInstalled() {
    return !!spawnSync('yarn', ['bin']).output
  }

  setIntegrationConfig = (config: { integrationTitle: string; orgId: string; integrationId: string }) => {
    const { integrationTitle, orgId, integrationId } = config
    if (this.rootPathRc) {
      fs.writeFileSync(this.rootPathRc, ini.stringify({ integrationTitle, orgId, integrationId }))
    }
  }

  storeBearerConfig = async (config: { Username?: string; ExpiresAt: number; authorization: any }) => {
    // deprecated
    const { Username, ExpiresAt, authorization } = config
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

  storeToken = async (token: TAccessToken) => {
    try {
      const file = this.bearerConfig.config || path.join(os.homedir(), '.bearerrc')
      let config = {}

      if (fs.existsSync(file)) {
        const configContent = await readFile(file, { encoding: 'utf8' })
        config = ini.parse(configContent)
      }
      const tokenWithExpires = { ...token, expires_at: Date.now() + token.expires_in * 1000 }
      await writeFile(file, ini.stringify({ ...config, token: tokenWithExpires }))
    } catch (e) {
      console.error('Error while writing token', e)
    }
  }

  async getToken(): Promise<TAccessToken | undefined> {
    const file = this.bearerConfig.config || path.join(os.homedir(), '.bearerrc')
    if (fs.existsSync(file)) {
      const configContent = await readFile(file, { encoding: 'utf8' })
      return ini.parse(configContent).token as TAccessToken
    }
    return undefined
  }
}

export default (runPath: string = process.cwd()): { constants: BaseConfig; config: Config } => {
  const { BEARER_ENV = 'production' } = process.env
  const config = new Config(runPath)
  return {
    config,
    constants: configs[BEARER_ENV as BearerEnv]
  }
}
