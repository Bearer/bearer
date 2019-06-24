import { captureHttps, setupFunctionIdentifiers } from '../src/index'
import { event } from './helpers/utils'

jest.mock('../src/constants')

describe('captureHttp', () => {
  it('overrides http and https modules', async () => {
    expect(require('http')._bearerLoading).toBeUndefined()
    expect(require('https')._bearerLoading).toBeUndefined()

    captureHttps()

    expect(require('http')._bearerLoading).toBeDefined()
    expect(require('https')._bearerLoading).toBeDefined()
  })
})

describe('setupFunctionIdentifiers', () => {
  it('setup all function identifiers', () => {
    const functionEvent = event()

    setupFunctionIdentifiers(functionEvent)

    expect(process.env.clientId).toEqual(functionEvent.context.clientId)
    expect(process.env.scenarioUuid).toEqual(functionEvent.context.integrationUuid)
  })
})
