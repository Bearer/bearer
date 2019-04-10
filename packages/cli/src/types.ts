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
  LoginDomain: string
}

export type BearerConfig = {
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
  integrationTitle: string
  open: boolean
  configs: string[]
  rootPathRc: string | null
  config: string
}

export type AuthConfig = {
  authType: Authentication
}

export type IntegrationBuildEnv = {
  BEARER_INTEGRATION_ID: string
  BEARER_INTEGRATION_HOST: string
  BEARER_AUTHORIZATION_HOST: string
  CDN_HOST: string
}

export type TAccessToken = {
  access_token: string
  id_token: string
  refresh_token: string
  scope: string
  expires_in: number
  expires_at: number
  token_type: string
}
