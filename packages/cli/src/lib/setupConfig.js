const findUp = require('find-up')
const rc = require('rc')
const os = require('os')

const fs = require('fs')
const ini = require('ini')
const path = require('path')

const scenarioConfig = rc('scenario')
const bearerConfig = rc('bearer')
const rootPathRc = findUp.sync('.scenariorc')

const IntegrationServiceHost = 'https://int.staging.bearer.sh/'

let setup = {
  DeploymentUrl: 'https://developer.staging.bearer.sh/v1/',
  IntegrationServiceHost,
  IntegrationServiceUrl: `${IntegrationServiceHost}api/v1/`,
  BearerEnv: 'staging',
  DeveloperPortalAPIUrl: 'https://app.bearer.sh/graphql'
}

module.exports = () => {
  const { BEARER_ENV } = process.env

  if (BEARER_ENV === 'dev') {
    setup = {
      DeploymentUrl: 'https://developer.dev.bearer.sh/v1/',
      IntegrationServiceUrl: 'https://int.dev.bearer.sh/api/v1/',
      IntegrationServiceHost: 'https://int.dev.bearer.sh/',
      DeveloperPortalAPIUrl: 'https://app.bearer.sh/graphql',
      BearerEnv: 'dev'
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
