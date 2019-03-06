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

export type IntegrationConfig = {
  integrationId: string
  integrationUuid: string
  orgId: string
  integrationTitle: string
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
  integrationConfig: IntegrationConfig
  orgId: string | undefined
  integrationTitle: string | undefined
  integrationId: string | undefined
  integrationUuid: string
  rootPathRc: string
  hasIntegrationLinked: boolean
  setIntegrationConfig(config: any): void
  storeBearerConfig(config: any): void
}
