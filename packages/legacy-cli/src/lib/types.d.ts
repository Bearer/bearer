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
  ExpiresAt: number
  authorization: {
    AuthenticationResult: any
  }
  open: false
  configs: Array<string>
  config: string
}

export type ScenarioConfig = {
  scenarioId: string
  scenarioUuid: string
  orgId: string
  scenarioTitle: string
  open: boolean
  configs: Array<string>
  rootPathRc: string
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
  rootPathRc: string
  hasScenarioLinked: boolean
  setScenarioConfig(config: any): void
  storeBearerConfig(config: any): void
}
