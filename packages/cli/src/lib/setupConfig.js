const findUp = require('find-up')
const rc = require('rc')
const os = require('os')

const fs = require('fs')
const ini = require('ini')
const path = require('path')

const scenarioConfig = rc('scenario')
const bearerConfig = rc('bearer')
const rootPathRc = findUp.sync('.scenariorc')

let setup = {
  DeploymentUrl: 'https://developer.staging.bearer.sh/v1/',
  IntegrationServiceHost: 'https://int.staging.bearer.sh/',
  IntegrationServiceUrl: 'https://int.staging.bearer.sh/api/v1/',
  DeveloperPortalAPIUrl: 'https://app.staging.bearer.sh/graphql',
  BearerEnv: 'staging'
}

module.exports = () => {
  const { BEARER_ENV } = process.env

  if (BEARER_ENV === 'dev') {
    setup = {
      DeploymentUrl: 'https://developer.dev.bearer.sh/v1/',
      IntegrationServiceHost: 'https://int.dev.bearer.sh/',
      IntegrationServiceUrl: 'https://int.dev.bearer.sh/api/v1/',
      DeveloperPortalAPIUrl: 'https://app.staging.bearer.sh/graphql',
      BearerEnv: 'dev'
    }
  }

  if (BEARER_ENV === 'production') {
    setup = {
      DeploymentUrl: 'https://developer.bearer.sh/v1/',
      IntegrationServiceHost: 'https://int.bearer.sh/',
      IntegrationServiceUrl: 'https://int.bearer.sh/api/v1/',
      DeveloperPortalAPIUrl: 'https://app.bearer.sh/graphql',
      BearerEnv: 'production'
    }
  }

  return {
    ...setup,
    HandlerBase: 'index.js',
    bearerConfig,
    scenarioConfig,
    rootPathRc,
    storeBearerConfig(config) {
      const { OrgId, Username, ExpiresAt, authorization } = config
      fs.writeFileSync(
        this.bearerConfig.config || path.join(os.homedir(), '.bearerrc'),
        ini.stringify({ OrgId, Username, ExpiresAt, authorization })
      )
    }
  }
}
