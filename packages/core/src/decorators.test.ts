import { BearerWindow } from '@bearer/types'

import Bearer from './bearer'
import { BearerFetch, _BackendFunction } from './decorators/functions'

declare const window: BearerWindow
declare const global: { fetch: any }

const INTEGRATION_ID = '1234'
const commonParams = {
  body: '{}',
  method: 'POST',
  credentials: 'include',
  headers: {
    'content-type': 'application/json',
    'user-agent': 'Bearer'
  }
}

describe('Function decorator', () => {
  beforeAll(() => {
    Bearer.init({ integrationHost: process.env.API_HOST })
    Bearer.instance.allowIntegrationRequests(true)
  })

  class FunctionDecorated {
    @_BackendFunction('getCollectionFunction')
    getCollectionFunctionProp: BearerFetch

    get INTEGRATION_ID(): string {
      return INTEGRATION_ID
    }

    get setupId() {
      return 'setup-id-from-props'
    }
  }

  const decoratedInstance: any = new FunctionDecorated()
  describe('FetchData', () => {
    const collection = [{ id: 42 }]

    beforeEach(() => {
      global.fetch.resetMocks()
      global.fetch.mockResponseOnce(JSON.stringify({ data: collection }))
    })

    it('adds a method', () => {
      expect(typeof decoratedInstance.getCollectionFunctionProp).toBe('function')
    })

    it('calling methods return a promise', () => {
      expect(decoratedInstance.getCollectionFunctionProp().constructor).toBe(Promise)
    })

    it('uses FetchData', async () => {
      const success = jest.fn()
      window.bearer = { clientId: '42', load: jest.fn() }

      await decoratedInstance
        .getCollectionFunctionProp({ page: 1 })
        .then(success)
        .catch(a => console.log(a))

      expect(global.fetch).toBeCalledWith(
        'https://localhost:5555/api/v3/functions/1234/getCollectionFunction?page=1&setupId=setup-id-from-props&clientId=42',
        commonParams
      )

      expect(success).toBeCalledWith({ data: collection, referenceId: null })
    })

    it('allows custom setupId', async () => {
      const success = jest.fn()
      window.bearer = { clientId: '42', load: jest.fn() }

      await decoratedInstance
        .getCollectionFunctionProp({ page: 1, setupId: 'custom-setupId' })
        .then(success)
        .catch(a => console.log(a))

      expect(global.fetch).toBeCalledWith(
        'https://localhost:5555/api/v3/functions/1234/getCollectionFunction?page=1&setupId=custom-setupId&clientId=42',
        commonParams
      )

      expect(success).toBeCalledWith({ data: collection, referenceId: null })
    })
  })
})
