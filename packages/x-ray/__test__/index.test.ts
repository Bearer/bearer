import { captureHttps, setupFunctionIdentifiers } from '../src/index'

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
    const event = {
      context: {
        clientId: 'myClientId',
        integrationUuid: 'myIntegrationUuid'
      }
    }

    setupFunctionIdentifiers(event)

    expect(process.env.clientId).toEqual('myClientId')
    expect(process.env.scenarioUuid).toEqual('myIntegrationUuid')
  })
})
