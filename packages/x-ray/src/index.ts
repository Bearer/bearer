import debug from '@bearer/logger'
import { overrideRequestMethod, overrideGetMethod } from './http-overrides'
const logger = debug('intents')

export const captureHttps = (module: any, event: any) => {
  logger(`Ìnject ${event.params.clientId} and ${event.params.scenarioUuid}`)
  process.env.clientId = event.params.clientId
  process.env.scenarioUuid = event.params.scenarioUuid

  if (module._request || module._get) {
    // This because the cold lambda start
    // Avoid overriding methods for next calls
    return
  }

  overrideRequestMethod(module)
  overrideGetMethod(module)
}
