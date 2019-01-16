import mockAxios from 'jest-mock-axios'

import { TLambdaEvent, TNONEAuthContext } from '../declaration'
import { RetrieveState } from './retrieve-state'

describe('Intents', () => {
  afterEach(() => {
    mockAxios.reset()
  })

  describe('RetrieveState', () => {
    it('export a class', () => {
      expect(RetrieveState).toBeTruthy()
    })
  })

  describe('existing data', () => {
    const referenceId = 'SPONGE_BOB'

    const event: TLambdaEvent<TNONEAuthContext> = {
      queryStringParameters: { referenceId },
      context: {
        bearerBaseURL: 'https://void.bearer.sh',
        signature: 'encrypted',
        authAccess: undefined
      }
    }

    const action = (_context: any, _params: any, state: any, callback: (state: any) => void) => {
      callback({ state, data: state })
    }

    it('returns meta + data', async done => {
      const lambdaCallback = jest.fn((a, b) => {
        expect(a).toBe(null)
        expect(b).toMatchObject({
          meta: { referenceId: 'SPONGE_BOB' },
          data: { title: 'Sponge bob' }
        })
        done()
      })

      RetrieveState.intent(action)(event, {}, lambdaCallback)
      expect(mockAxios.get).toHaveBeenCalledWith('api/v2/items/SPONGE_BOB', { params: { signature: 'encrypted' } })
      mockAxios.mockResponse({ data: { Item: { referenceId, title: 'Sponge bob' } } })
    })

    it('not found retuns error', async done => {
      const lambdaCallback = jest.fn((a, b) => {
        expect(a).toBe(null)
        expect(b).toMatchObject({
          statusCode: 404
        })
        done()
      })

      RetrieveState.intent(action)(event, {}, lambdaCallback)
      expect(mockAxios.get).toHaveBeenCalledWith('api/v2/items/SPONGE_BOB', { params: { signature: 'encrypted' } })
      mockAxios.mockError({ data: { Item: null }, status: 404 })
    })
  })
})
