import Bearer from './Bearer'
import { Intent, BearerFetch, IntentType } from './decorators/Intent'
import { BearerComponent } from './decorators/Component'
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

  @BearerComponent
  class IntentDecorated {
    @Intent('getCollectionIntent') getCollectionIntentProp: BearerFetch
    @Intent('getResourceIntent', IntentType.GetResource)
    getResourceIntentProp: BearerFetch

    SCENARIO_ID: string = SCENARIO_ID
  }

  const decoratedInstance = new IntentDecorated()

  describe('GetCollectionIntent', () => {
    const collection = [{ id: 42 }]

    beforeEach(() => {
      fetch.resetMocks()
      fetch.mockResponseOnce(JSON.stringify(collection))
    })

    it('adds a method', () => {
      expect(typeof decoratedInstance.getCollectionIntentProp).toBe('function')
    })

    it('calling methods return a promise', () => {
      expect(decoratedInstance.getCollectionIntentProp().constructor).toBe(
        Promise
      )
    })

    it('uses GetResourceIntent', async () => {
      const success = jest.fn()

      await decoratedInstance
        .getCollectionIntentProp({ page: 1 })
        .then(success)
        .catch(a => console.log(a))

      expect(fetch).toBeCalledWith(
        'http://localhost:5555/api/v1/BEARER_SCENARIO_ID/getCollectionIntent?page=1&integrationId=42',
        commonParams
      )

      expect(success).toBeCalledWith({ items: collection })
    })
  })

  describe('GetResourceIntent', () => {
    const item = { id: 42 }

    beforeEach(() => {
      fetch.resetMocks()
      fetch.mockResponseOnce(JSON.stringify(item))
    })

    it('adds a method', () => {
      expect(typeof decoratedInstance.getResourceIntentProp).toBe('function')
    })

    it('calling methods return a promise', () => {
      expect(decoratedInstance.getResourceIntentProp().constructor).toBe(
        Promise
      )
    })

    it('uses GetResourceIntent', async () => {
      const success = jest.fn()

      await decoratedInstance
        .getResourceIntentProp()
        .then(success)
        .catch(a => console.log(a))

      expect(fetch).toBeCalledWith(
        'http://localhost:5555/api/v1/BEARER_SCENARIO_ID/getResourceIntent?integrationId=42',

        commonParams
      )

      expect(success).toBeCalledWith({ object: item })
    })
  })

  describe('BearerComponent', () => {
    it('rejects if no scenario_id is not provided', async () => {
      class MissingBearerComponentDecorator {
        @Intent('rejectedIntent') rejectedIntent: BearerFetch
        SCENARIO_ID: string = null
      }

      const instance = new MissingBearerComponentDecorator()
      await expect(instance.rejectedIntent()).rejects.toThrow(
        'Scenario ID is missing. Add BearerComponent above @Component({...}) decorator'
      )
    })
  })
})
