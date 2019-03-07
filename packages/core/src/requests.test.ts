import { BearerWindow } from '@bearer/types'
// import fetch from 'jest-fetch-mock'

import Bearer from './bearer'
import { functionRequest, itemRequest } from './requests'

const functionName = 'anFunction'
const integrationId = 'aIntegrationId'
const setupId = '1234'
declare const window: BearerWindow & { fetch: any }
declare const global: { fetch: any }

describe('requests', () => {
  beforeEach(() => {
    global.fetch.resetMocks()
    Bearer.init({ integrationHost: process.env.API_HOST, secured: true })
    Bearer.instance.allowIntegrationRequests(true)
  })

  describe('itemRequest', () => {
    it('returns a function', () => {
      const aRequest = itemRequest()

      expect(typeof aRequest).toBe('function')
    })

    it('calls host + functionName + params', async () => {
      const aRequest = itemRequest()
      global.fetch.mockResponseOnce(JSON.stringify({}))
      window.bearer = { clientId: '42', load: jest.fn() }

      await aRequest({}, {})

      expect(global.fetch).toBeCalledWith('https://localhost:5555/api/v1/items?clientId=42&secured=true', {
        credentials: 'include',
        headers: {
          'content-type': 'application/json',
          'user-agent': 'Bearer'
        }
      })
    })
  })

  describe('functionRequest', () => {
    it('returns a function', () => {
      const aRequest = functionRequest({ functionName, integrationId, setupId })

      expect(typeof aRequest).toBe('function')
    })

    it('calls host + functionName + params', async () => {
      const aRequest = functionRequest({ functionName, integrationId, setupId })
      global.fetch.mockResponseOnce(JSON.stringify({}))
      window.bearer = { clientId: '42', load: jest.fn() }

      await aRequest({ page: 1 }, {})

      expect(global.fetch).toBeCalledWith(
        // tslint:disable-next-line:max-line-length
        'https://localhost:5555/api/v3/functions/aIntegrationId/anFunction?page=1&setupId=1234&clientId=42&secured=true',
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
