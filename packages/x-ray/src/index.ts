import logger from './logger'
import { overrideRequestMethod, overrideGetMethod } from './http-overrides'

export const captureHttps = (module: any, event: any) => {
  const context = event.context || {}
  const { clientId, integrationUuid } = context
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
