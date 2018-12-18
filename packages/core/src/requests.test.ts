import { BearerWindow } from '@bearer/types'
// import fetch from 'jest-fetch-mock'

import Bearer from './bearer'
import { intentRequest } from './requests'

const intentName = 'anIntent'
const scenarioId = 'aScenarioId'
const setupId = '1234'
declare const window: BearerWindow & { fetch: any }
declare const global: { fetch: any }

describe('requests', () => {
  beforeEach(() => {
    global.fetch.resetMocks()
    Bearer.init({ integrationHost: process.env.API_HOST, secured: true })
    Bearer.instance.allowIntegrationRequests(true)
  })

  describe('intentRequest', () => {
    it('returns a function', () => {
      const aRequest = intentRequest({ intentName, scenarioId, setupId })

      expect(typeof aRequest).toBe('function')
    })

    it('calls host + intentName + params', async () => {
      const aRequest = intentRequest({ intentName, scenarioId, setupId })
      global.fetch.mockResponseOnce(JSON.stringify({}))
      window.bearer = { clientId: '42', load: jest.fn() }

      await aRequest({ page: 1 }, {})

      expect(global.fetch).toBeCalledWith(
        'https://localhost:5555/api/v1/aScenarioId/anIntent?page=1&setupId=1234&clientId=42&secured=true',
        {
          credentials: 'include',
          headers: {
            'content-type': 'application/json',
            'user-agent': 'Bearer'
          }
        }
      )
    })
  })
})
