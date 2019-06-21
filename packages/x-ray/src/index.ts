import logger from './logger'
import { overrideRequestMethod } from './http-overrides'

let check: number

export const captureHttps = () => {
  const httpModule = require('http')
  const httpsModule = require('https')
  if (httpModule._bearerLoading === check && check) {
    logger('%j', { message: 'http module has already been loaded', application: 'x-ray' })
    return
  }
  httpModule._bearerLoading = check = Math.random()
  logger('%j', { message: 'Overriding request and get methods', application: 'x-ray' })
  overrideRequestMethod(httpModule)
  overrideRequestMethod(httpsModule)
}

export const setupFunctionIdentifiers = function(event: any) {
  const context = event.context || {}
  const { clientId, integrationUuid } = context
  logger('%j', { message: `Ìnject ${clientId} and ${integrationUuid}`, application: 'x-ray' })
  process.env.clientId = clientId
  process.env.scenarioUuid = integrationUuid
}

export default captureHttps
