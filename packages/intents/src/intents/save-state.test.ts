import mockAxios from 'jest-mock-axios'

import { TLambdaEvent, TAPIKEYAuthContext } from '../declaration'
import { SaveState } from './save-state'

describe('Intents', () => {
  describe('SaveState', () => {
    it('export a class', () => {
      expect(SaveState).toBeTruthy()
    })
  })

  describe('.intent', () => {
    const referenceId = 'SPONGE_BOB'

    const action = (_context: any, _params: any, body: any, state: any, callback: (state: any) => void) => {
      callback({ state, data: { ...state, ...body } })
    }

    const event: TLambdaEvent<TAPIKEYAuthContext> = {
      queryStringParameters: { referenceId },
      context: {
        bearerBaseURL: 'https://void.bearer.sh',
        signature: 'encrypted',
        authAccess: {
          apiKey: 'none'
        }
      },
      body: {
        pullRequests: []
      }
    }

    describe('reference exists', () => {
      const savedData = {
        data: 'server says hello!'
      }

      function mockRetrieveSucces() {
        expect(mockAxios.get).toHaveBeenCalledWith('api/v2/items/SPONGE_BOB', { params: { signature: 'encrypted' } })
        const responseObj = { data: { Item: savedData } }
        mockAxios.mockResponse(responseObj, mockAxios.lastReqGet())
      }

      it.skip('retrieves reference, update it and return payload', async done => {
        const lambdaCallback = jest.fn((a, b) => {
          console.log('ok', a, b)
          done()
        })
        // expect.assertions(4)

        SaveState.intent(action)(event, {}, lambdaCallback)
        mockRetrieveSucces()
        // expect(mockAxios.put).toHaveBeenCalled()
        expect(mockAxios.put).toHaveBeenCalledWith(
          'api/v2/items/SPONGE_BOB',
          {
            //   ReadAllowed: true,
            //   ...savedData,
            //   pullRequests: []
          },
          { params: { signature: 'encrypted' } }
        )
        mockAxios.mockResponse({}) // data is unused

        expect(lambdaCallback).toHaveBeenCalledWith(null, {
          meta: { referenceId },
          data: { ...savedData, pullRequests: [] }
        })
        // expect(lambdaCallback).toHaveBeenCalled()
      })

      it.skip('retrieves reference, fails update and return error', () => {
        const lambdaCallback = jest.fn()

        SaveState.intent(action)(event, {}, lambdaCallback)

        mockRetrieveSucces()

        expect(mockAxios.put).toHaveBeenCalledWith(
          'api/v2/items/SPONGE_BOB',
          {
            ReadAllowed: true,
            ...savedData,
            pullRequests: []
          },
          { params: { signature: 'encrypted' } }
        )
        mockAxios.mockError({ status: 500 })

        expect(lambdaCallback).toHaveBeenCalledWith(expect.any(String), { error: 'Cannot update data' })
      })

      it.skip('retrieves data, action throw error, returns error', () => {
        const lambdaCallback = jest.fn()
        const failingAction = () => {
          throw new Error('Error')
        }

        SaveState.intent(failingAction)(event, {}, lambdaCallback)

        mockRetrieveSucces()

        expect(mockAxios.put).not.toHaveBeenCalled()
        expect(lambdaCallback).toHaveBeenCalledWith(expect.any(String), { error: 'Intent action is failing' })
      })
    })

    describe('reference does not exist', () => {
      function mockRetrieveFailure() {
        expect(mockAxios.get).toHaveBeenCalledWith('api/v2/items/SPONGE_BOB', { params: { signature: 'encrypted' } })
        const responseObj = { data: { Item: null }, status: 404 }
        mockAxios.mockError(responseObj)
      }

      it.skip('create reference and return payload', () => {
        const lambdaCallback = jest.fn()

        SaveState.intent(action)(event, {}, lambdaCallback)

        // mockRetrieveFailure()
        expect(mockAxios.get).toHaveBeenCalledWith('api/v2/items/SPONGE_BOB', mockAxios.lastReqGet())
        expect(mockAxios.post).toHaveBeenCalledWith(
          'api/v2/items',
          { ReadAllowed: true, pullRequests: [] },
          { params: { signature: 'encrypted' } }
        )

        const responseObj = { data: { Item: null }, status: 404 }
        mockAxios.mockError(responseObj)
        mockAxios.mockResponse({ data: { Item: { referenceId: 'PATRICK' } } })

        expect(mockAxios.post).toHaveBeenCalledWith()
        expect(lambdaCallback).toHaveBeenCalledWith(null, { meta: { referenceId: 'PATRICK', data: { ...event.body } } })
      })

      it.skip('fails create reference and return error', () => {
        const lambdaCallback = jest.fn()

        SaveState.intent(action)(event, {}, lambdaCallback)

        mockRetrieveFailure()

        expect(mockAxios.post).toHaveBeenCalledWith(
          'api/v2/items',
          { ReadAllowed: true, pullRequests: [] },
          { params: { signature: 'encrypted' } }
        )
        const responseObj = {}
        mockAxios.mockError(responseObj)
        expect(lambdaCallback).toHaveBeenCalledWith(expect.any(String), { error: 'Cannot save data' })
      })
      it.skip('action throw error, returns error', () => {
        const lambdaCallback = jest.fn()
        const failingAction = () => {
          throw new Error('Error')
        }

        SaveState.intent(failingAction)(event, {}, lambdaCallback)
        mockRetrieveFailure()

        expect(mockAxios.post).not.toHaveBeenCalled()
        expect(lambdaCallback).toHaveBeenCalledWith(expect.any(String), { error: 'Intent action is failing' })
      })
    })

    describe('fails userData service call', () => {
      it.skip('returns error', () => {
        const lambdaCallback = jest.fn()

        SaveState.intent(action)(event, {}, lambdaCallback)
        expect(mockAxios.get).toHaveBeenCalledWith('api/v2/items/SPONGE_BOB', { params: { signature: 'encrypted' } })
        const responseObj = { data: { Item: null }, status: 500 }
        mockAxios.mockError(responseObj)

        expect(mockAxios.post).not.toHaveBeenCalledWith()
        expect(mockAxios.put).not.toHaveBeenCalledWith()
        expect(lambdaCallback).toHaveBeenCalledWith(null, { error: 'Error whild trying to fetch current state' })
      })
    })
  })
})
