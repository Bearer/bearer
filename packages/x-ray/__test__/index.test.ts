import { bearerOverride, setupFunctionIdentifiers } from '../src/index'

jest.mock('../src/constants')

describe('bearerOverride', () => {
  it('overrides modules', async () => {
    expect(require('http')._bearerLoading).toBeUndefined()
    expect(require('https')._bearerLoading).toBeUndefined()
    expect(require('console')._bearerLoading).toBeUndefined()

    bearerOverride()

    expect(require('http')._bearerLoading).toBeDefined()
    expect(require('console')._bearerLoading).toBeDefined()
  })
})

describe('setupFunctionIdentifiers', () => {
  it('setup all function identifiers', () => {
    const event = {
      context: {
        organizationIdentifier: 'organizationIdentifier',
        clientId: 'myClientId',
        integrationUuid: 'myIntegrationUuid'
      }
    }

    setupFunctionIdentifiers(event)

    expect(process.env.clientId).toEqual('organizationIdentifier')
    expect(process.env.scenarioUuid).toEqual('myIntegrationUuid')
  })

  describe('when organizationIdentifier is missing in context', () => {
    it('setup fallbacks to clientId', () => {
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
})
