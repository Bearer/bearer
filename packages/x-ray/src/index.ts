import logger from './logger'
import { overrideGetMethod, overrideRequestMethod } from './http-overrides'

export const captureHttps = function(module: any) {
  if (module._request && module._get) {
    logger('%j', { message: 'request & get already overridden', application: 'x-ray' })
    // This because the cold lambda start
    // Avoid overriding methods for next calls
    return
  }

  logger('%j', { message: 'Override request and get methods', application: 'x-ray' })
  overrideRequestMethod(module)
  overrideGetMethod(module)
}

export const setupFunctionIdentifiers = function(event: any) {
  const context = event.context || {}
  const { clientId, integrationUuid } = context
  logger('%j', { message: `Ìnject ${clientId} and ${integrationUuid}`, application: 'x-ray' })
  process.env.clientId = clientId
  process.env.scenarioUuid = integrationUuid
}
