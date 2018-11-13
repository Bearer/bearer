import { BearerWindow } from '@bearer/types'

import Bearer from './bearer'
import { BearerFetch, Intent } from './decorators/Intent'

declare const window: BearerWindow
declare const global: { fetch: any }

const SCENARIO_ID = '1234'
const commonParams = {
  body: '{}',
  method: 'POST',
  credentials: 'include',
  headers: {
    'content-type': 'application/json',
    'user-agent': 'Bearer'
  }
}

describe('Intent decorator', () => {
  beforeAll(() => {
    Bearer.init({ integrationHost: process.env.API_HOST })
    Bearer.instance.allowIntegrationRequests(true)
  })

  class IntentDecorated {
    @Intent('getCollectionIntent')
    getCollectionIntentProp: BearerFetch

    get SCENARIO_ID(): string {
      return SCENARIO_ID
    }
  }

  const decoratedInstance = new IntentDecorated()

  describe('FetchData', () => {
    const collection = [{ id: 42 }]

    beforeEach(() => {
      global.fetch.resetMocks()
      global.fetch.mockResponseOnce(JSON.stringify({ data: collection }))
    })

    it('adds a method', () => {
      expect(typeof decoratedInstance.getCollectionIntentProp).toBe('function')
    })

    it('calling methods return a promise', () => {
      expect(decoratedInstance.getCollectionIntentProp().constructor).toBe(Promise)
    })

    it('uses FetchData', async () => {
      const success = jest.fn()
      window.bearer = { clientId: '42', load: jest.fn() }

      await decoratedInstance
        .getCollectionIntentProp({ page: 1 })
        .then(success)
        .catch(a => console.log(a))

      expect(global.fetch).toBeCalledWith(
        'https://localhost:5555/api/v1/1234/getCollectionIntent?page=1&setupId=&clientId=42',
        commonParams
      )

      expect(success).toBeCalledWith({ data: collection, referenceId: null })
    })
  })
})
