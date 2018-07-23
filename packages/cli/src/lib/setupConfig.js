const findUp = require('find-up')
const rc = require('rc')
const os = require('os')

const fs = require('fs')
const ini = require('ini')
const path = require('path')

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
  const configuration = {
    ...setup,
    HandlerBase: 'index.js',
    get OrgId() {
      return this.bearerConfig.OrgId
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
    get rootPathRc() {
      return findUp.sync('.scenariorc')
    },
    get credentials() {
      const { Username, infrastructurePassword } = this.bearerConfig
      return { Username, infrastructurePassword }
    },
    storeBearerConfig(config) {
      const {
        OrgId,
        Username,
        ExpiresAt,
        authorization,
        infrastructurePassword = ''
      } = config
      fs.writeFileSync(
        this.bearerConfig.config || path.join(os.homedir(), '.bearerrc'),
        ini.stringify({
          OrgId,
          Username,
          ExpiresAt,
          authorization,
          infrastructurePassword
        })
      )
    },
    setScenarioConfig(config) {
      const { scenarioTitle, scenarioId } = config
      console.log(this.rootPathRc)
      fs.writeFileSync(
        this.rootPathRc,
        ini.stringify({ scenarioTitle, scenarioId })
      )
    }
  }
  return configuration
}
