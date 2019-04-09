import bearer from '@bearer/js'
import { BearerFetch, _BackendFunction } from './decorators/functions'

declare const global: { fetch: any }

const INTEGRATION_ID = '1234'

describe('Function decorator', () => {
  beforeAll(() => {
    bearer('client-id', { integrationHost: 'HOST' })
    // @ts-ignore
    window.bearer = bearer
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

      await decoratedInstance.getCollectionFunctionProp({ page: 1 }).then(success)

      expect(global.fetch).toBeCalledWith(
        'HOST/api/v4/functions/1234/getCollectionFunction?setupId=setup-id-from-props&clientId=client-id',
        {
          body: '{"page":1}',
          credentials: 'include',
          headers: { 'content-type': 'application/json' },
          method: 'POST'
        }
      )

      expect(success).toBeCalledWith({ data: collection })
    })

    it('allows custom setupId', async () => {
      const success = jest.fn()

      await decoratedInstance.getCollectionFunctionProp({ page: 1, setupId: 'custom-setupId' }).then(success)

      expect(global.fetch).toBeCalledWith(
        'HOST/api/v4/functions/1234/getCollectionFunction?setupId=setup-id-from-props&clientId=client-id',
        {
          body: '{"page":1,"setupId":"custom-setupId"}',
          credentials: 'include',
          headers: { 'content-type': 'application/json' },
          method: 'POST'
        }
      )

      expect(success).toBeCalledWith({ data: collection })
    })
  })

  describe('when error returned', () => {
    beforeAll(() => {
      global.fetch.resetMocks()
      global.fetch.mockResponseOnce(JSON.stringify({ error: 'error from api' }))
    })

    it('pass through catch', async () => {
      const error = jest.fn()

      await decoratedInstance.getCollectionFunctionProp({ page: 1 }).catch(error)

      expect(global.fetch).toBeCalledWith(
        'HOST/api/v4/functions/1234/getCollectionFunction?setupId=setup-id-from-props&clientId=client-id',
        {
          body: '{"page":1}',
          credentials: 'include',
          headers: { 'content-type': 'application/json' },
          method: 'POST'
        }
      )

      expect(error).toBeCalledWith({ error: { error: 'error from api' } })
    })
  })
})
