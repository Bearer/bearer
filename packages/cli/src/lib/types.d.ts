export type ScenarioConfig = {
  DeploymentUrl: string
  IntegrationServiceHost: string
  IntegrationServiceUrl: string
  BearerEnv: string
  DeveloperPortalAPIUrl: string
  HandlerBase: string
  bearerConfig: {
    OrgId: string
    Username: string
    ExpiresAt: string
    authorization: {
      AuthenticationResult: any
    }
    open: false
    configs: Array<string>
    config: string
  }
  scenarioConfig: {
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
}
