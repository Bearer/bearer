export const enum AuthType {
  OAuth2 = 'oauth2',
  Basic = 'basicauth',
  ApiKey = 'apikey',
  NoAuth = 'noauth'
}

export type BearerEnv = 'dev' | 'staging' | 'production'

export type BaseConfig = {
  DeploymentUrl: string
  IntegrationServiceHost: string
  IntegrationServiceUrl: string
  BearerEnv: string
  DeveloperPortalAPIUrl: string
  DeveloperPortalUrl: string
  CdnHost: string
}

export type BearerConfig = {
  OrgId: string
  Username: string
  ExpiresAt: string
  authorization: {
    AuthenticationResult: any
  }
  open: false
  configs: Array<string>
  config: string
}

export type ScenarioConfig = {
  scenarioId: string
  scenarioUuid: string | null
  orgId: string
  scenarioTitle: string
  open: boolean
  configs: Array<string>
  rootPathRc: string | null
  storeBearerConfig: any
  config: string
}

export type Config = BaseConfig & {
  isYarnInstalled: boolean
  command: 'yarn' | 'npm'
  bearerConfig: BearerConfig
  scenarioConfig: ScenarioConfig
  orgId: string | undefined
  scenarioTitle: string | undefined
  scenarioId: string | undefined
  scenarioUuid: string
  rootPathRc: string | null
  hasScenarioLinked: boolean
  setScenarioConfig(config: any): void
  storeBearerConfig(config: any): void
}
