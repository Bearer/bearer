const findUp = require('find-up')
const rc = require('rc')
const os = require('os')

const fs = require('fs')
const ini = require('ini')
const path = require('path')

let setup = {
  DeploymentUrl: 'https://developer.staging.bearer.sh/v1/',
  IntegrationServiceHost: 'https://int.staging.bearer.sh/',
  IntegrationServiceUrl: 'https://int.staging.bearer.sh/api/v1/',
  DeveloperPortalAPIUrl: 'https://app.staging.bearer.sh/graphql',
  DeveloperPortalUrl: 'https://app.staging.bearer.sh/',
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
      DeveloperPortalUrl: 'https://app.staging.bearer.sh/',
      BearerEnv: 'dev'
    }
  }

  if (BEARER_ENV === 'production') {
    setup = {
      DeploymentUrl: 'https://developer.bearer.sh/v1/',
      IntegrationServiceHost: 'https://int.bearer.sh/',
      IntegrationServiceUrl: 'https://int.bearer.sh/api/v1/',
      DeveloperPortalAPIUrl: 'https://app.bearer.sh/graphql',
      DeveloperPortalUrl: 'https://app.bearer.sh/',
      BearerEnv: 'production'
    }
  }

  return {
    ...setup,
    HandlerBase: 'index.js',
    get orgId() {
      return this.scenarioConfig.orgId
    },
    get bearerConfig() {
      return rc('bearer')
    },
    get scenarioConfig() {
      return rc('scenario')
    },
    get scenarioTitle() {
      return this.scenarioConfig.scenarioTitle
    },
    get scenarioId() {
      return this.scenarioConfig.scenarioId
    },
    get scenarioUuid() {
      return `${this.orgId}-${this.scenarioId}`
    },
    get rootPathRc() {
      return findUp.sync('.scenariorc')
    },
    get credentials() {
      const { Username, infrastructurePassword } = this.bearerConfig
      return { Username, infrastructurePassword }
    },
    storeBearerConfig(config) {
      const { Username, ExpiresAt, authorization } = config //, infrastructurePassword = '' } = config
      fs.writeFileSync(
        this.bearerConfig.config || path.join(os.homedir(), '.bearerrc'),
        ini.stringify({
          Username,
          ExpiresAt,
          authorization
          // infrastructurePassword
        })
      )
    },
    setScenarioConfig(config) {
      const { scenarioTitle, orgId, scenarioId } = config
      fs.writeFileSync(this.rootPathRc, ini.stringify({ scenarioTitle, orgId, scenarioId }))
    }
  }
}
