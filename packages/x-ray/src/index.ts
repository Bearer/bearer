import logger from './logger'
import { overrideRequestMethod } from './http-overrides'
import { overrideConsole } from './console-overrides'

let check: string

export const bearerOverride = () => {
  const httpModule = require('http')
  // const httpsModule = require('https')
  const consoleModule = require('console')

  if (httpModule._bearerLoading === check && consoleModule._bearerLoading && check) {
    logger('%j', { message: 'modules have already been loaded', application: 'x-ray' })
    return
  }

  httpModule._bearerLoading = consoleModule._bearerLoading = check = [Math.random(), Date.now()].join('/')
  logger('%j', { message: 'Overriding request and get methods', application: 'x-ray' })

  overrideRequestMethod(httpModule)
  // AWS does the stuff for us
  // overrideRequestMethod(httpsModule)
  overrideConsole(consoleModule)
}

export const setupFunctionIdentifiers = function(event: any) {
  const context = event.context || {}
  const { clientId, integrationUuid, organizationIdentifier } = context
  logger('%j', {
    message: `Injecting ${JSON.stringify({ clientId, integrationUuid, organizationIdentifier })} `,
    application: 'x-ray'
  })
  process.env.clientId = organizationIdentifier || clientId
  process.env.scenarioUuid = integrationUuid
}

export default bearerOverride
