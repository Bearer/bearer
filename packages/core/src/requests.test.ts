import fetch from 'jest-fetch-mock'

import Bearer from './Bearer'
import { intentRequest } from './requests'

const intentName = 'anIntent'
const scenarioId = 'aScenarioId'
const setupId = '1234'

describe('requests', () => {
  beforeEach(() => {
    fetch.resetMocks()
    Bearer.init({ integrationHost: process.env.API_HOST, integrationId: '42' })
    Bearer.instance.allowIntegrationRequests()
  })

  describe('intentRequest', () => {
    it('returns a function', () => {
      const aRequest = intentRequest({ intentName, scenarioId, setupId })

      expect(typeof aRequest).toBe('function')
    })

    it('calls host + intentName + params', async () => {
      const aRequest = intentRequest({ intentName, scenarioId, setupId })
      fetch.mockResponseOnce(JSON.stringify({}))

      await aRequest({ page: 1 }, {})

      expect(window.fetch).toBeCalledWith(
        'https://localhost:5555/api/v1/aScenarioId/anIntent?page=1&setupId=1234&integrationId=42',
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
