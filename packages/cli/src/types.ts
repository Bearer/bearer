import Authentication from '@bearer/types/lib/authentications'

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
    AuthenticationResult?: {
      IdToken: string
      RefreshToken: string
      TokenType: string
      ExpiresIn: number
      AccessToken: string
    }
  }
  open: false
  configs: string[]
  config: string
}

export type IntegrationConfig = {
  integrationId: string
  integrationUuid: string | null
  orgId: string
  integrationTitle: string
  open: boolean
  configs: string[]
  rootPathRc: string | null
  storeBearerConfig: any
  config: string
}

export type Config = BaseConfig & {
  runPath: string
  isYarnInstalled: boolean
  isIntegrationLocation: boolean
  command: 'yarn' | 'npm'
  bearerConfig: BearerConfig
  integrationConfig: IntegrationConfig
  orgId: string | undefined
  integrationTitle: string | undefined
  integrationId: string | undefined
  integrationUuid: string
  rootPathRc: string | null
  hasIntegrationLinked: boolean
  setIntegrationConfig(config: any): void
  storeBearerConfig(config: any): void
}

export type AuthConfig = {
  authType: Authentication
  setupViews?: any[]
}

export type IntegrationBuildEnv = {
  BEARER_INTEGRATION_ID: string
  BEARER_INTEGRATION_TAG_NAME: string
  BEARER_INTEGRATION_HOST: string
  BEARER_AUTHORIZATION_HOST: string
  CDN_HOST: string
}
