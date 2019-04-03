import { spawnSync } from 'child_process'
import * as findUp from 'find-up'
import * as fs from 'fs'
import * as ini from 'ini'
import * as os from 'os'
import * as path from 'path'
import * as rc from 'rc'
import { promisify } from 'util'

import { BaseConfig, BearerConfig, BearerEnv, IntegrationConfig, TAccessToken } from '../types'
import { BEARER_ENV, CONFIGS } from './constants'

const writeFile = promisify(fs.writeFile)
const readFile = promisify(fs.readFile)

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
    return this.rootPath !== null
  }

  get bearerConfig(): BearerConfig {
    return rc('bearer')
  }

  get integrationConfig(): IntegrationConfig {
    return rc('integration', { config: this.integrationRc })
  }

  get integrationRc() {
    return path.join(this.integrationLocation, '.integrationrc')
  }

  get orgId(): string | undefined {
    return process.env.BEARER_ORG_ID || this.integrationConfig.orgId
  }

  get integrationTitle(): string | undefined {
    return this.integrationConfig.integrationTitle
  }

  get integrationId(): string | undefined {
    return process.env.BEARER_INTEGRATION_ID || this.integrationConfig.integrationId
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

  get rootPath(): string | null {
    return findUp.sync('.bearer', { cwd: this.integrationLocation })
  }

  get isYarnInstalled() {
    return !!spawnSync('yarn', ['bin']).output
  }

  setIntegrationConfig = (config: { integrationTitle: string; orgId: string; integrationId: string }) => {
    const { integrationTitle, orgId, integrationId } = config
    if (this.rootPath) {
      fs.writeFileSync(this.integrationRc, ini.stringify({ integrationTitle, orgId, integrationId }))
    }
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
  return {
    config: new Config(runPath),
    constants: CONFIGS[BEARER_ENV as BearerEnv]
  }
}
