import Bearer from './Bearer'
import { Intent, BearerFetch, IntentType } from './decorators/Intent'
import fetch from 'jest-fetch-mock'

const SCENARIO_ID = '1234'
const commonParams = {
  credentials: 'include',
  headers: {
    'content-type': 'application/json',
    'user-agent': 'Bearer'
  }
}

describe('Intent decorator', () => {
  beforeAll(() => {
    Bearer.init({ integrationHost: process.env.API_HOST, integrationId: '42' })
    Bearer.instance['allowIntegrationRequests']()
  })

  class IntentDecorated {
    @Intent('getCollectionIntent') getCollectionIntentProp: BearerFetch
    @Intent('getResourceIntent', IntentType.GetResource)
    getResourceIntentProp: BearerFetch

    get SCENARIO_ID(): string {
      return SCENARIO_ID
    }
  }

  const decoratedInstance = new IntentDecorated()

  describe('GetCollectionIntent', () => {
    const collection = [{ id: 42 }]

    beforeEach(() => {
      fetch.resetMocks()
      fetch.mockResponseOnce(JSON.stringify({ data: collection }))
    })

    it('adds a method', () => {
      expect(typeof decoratedInstance.getCollectionIntentProp).toBe('function')
    })

    it('calling methods return a promise', () => {
      expect(decoratedInstance.getCollectionIntentProp().constructor).toBe(Promise)
    })

    it('uses GetCollectionIntent', async () => {
      const success = jest.fn()

      await decoratedInstance
        .getCollectionIntentProp({ page: 1 })
        .then(success)
        .catch(a => console.log(a))

      expect(fetch).toBeCalledWith(
        'http://localhost:5555/api/v1/1234/getCollectionIntent?page=1&setupId=&integrationId=42',
        commonParams
      )

      expect(success).toBeCalledWith({ items: collection, referenceId: null })
    })
  })

  describe('GetResourceIntent', () => {
    const item = { id: 42 }

    beforeEach(() => {
      fetch.resetMocks()
      fetch.mockResponseOnce(JSON.stringify({ data: item }))
    })

    it('adds a method', () => {
      expect(typeof decoratedInstance.getResourceIntentProp).toBe('function')
    })

    it('calling methods return a promise', () => {
      expect(decoratedInstance.getResourceIntentProp().constructor).toBe(Promise)
    })

    it('uses GetResourceIntent', async () => {
      const success = jest.fn()

      await decoratedInstance
        .getResourceIntentProp()
        .then(success)
        .catch(a => console.log(a))

      expect(fetch).toBeCalledWith(
        'http://localhost:5555/api/v1/1234/getResourceIntent?setupId=&integrationId=42',

        commonParams
      )

      expect(success).toBeCalledWith({ object: item, referenceId: null })
    })
  })
})
