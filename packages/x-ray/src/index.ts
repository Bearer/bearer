import { overrideRequestMethod, overrideGetMethod } from './http-overrides'

export const captureHttps = (module: any, event: any) => {
  if (module._request || module._get) {
    // This because the cold lambda start
    // Avoid overriding methods for next calls
    process.env.clientId = event.params.clientId
    process.env.scenarioUuid = event.params.scenarioUuid
    return
  }

  overrideRequestMethod(module, event)
  overrideGetMethod(module, event)
}
