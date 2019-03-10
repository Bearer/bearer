import logger from './logger'
import { overrideRequestMethod, overrideGetMethod } from './http-overrides'

export const captureHttps = (module: any, event: any) => {
  const clientId = (event.params || {}).clientId
  const integrationUuid = (event.params || {}).scenarioUuid
  logger(`Ìnject ${clientId} and ${integrationUuid}`)
  process.env.clientId = clientId
  process.env.scenarioUuid = integrationUuid

  if (module._request && module._get) {
    // This because the cold lambda start
    // Avoid overriding methods for next calls
    return
  }

  overrideRequestMethod(module)
  overrideGetMethod(module)
}
