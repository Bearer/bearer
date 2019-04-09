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
  integrationTitle: string | undefined
  buid: string
  rootPathRc: string
}
