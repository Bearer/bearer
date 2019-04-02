import { BearerEnv, BaseConfig } from '../types'

export const LOGIN_CLIENT_ID = process.env.BEARER_LOGIN_CLIENT_ID || 'Wgll39KqWnJWud473wq7hZhiXxeNjEU7'
export const BEARER_ENV = process.env.BEARER_ENV || 'production'

export const CONFIGS: Record<BearerEnv, BaseConfig> = {
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
