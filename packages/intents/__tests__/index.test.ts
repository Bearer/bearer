import mockAxios from 'jest-mock-axios'

import { TLambdaEvent } from '../src/declaration'
import { RetrieveState, SaveState } from '../src/index'

afterEach(() => {
  mockAxios.reset()
})

describe('index', () => {
  it.skip('export SaveState', () => {
    expect(SaveState).toBeTruthy()
  })

  it.skip('export RetrieveState', () => {
    expect(RetrieveState).toBeTruthy()
  })
})

describe('SaveState', () => {
  describe('.intent', () => {
    const referenceId = 'SPONGE_BOB'

    const action = (_context: any, _params: any, body: any, state: any, callback: (state: any) => void) => {
      callback({ state, data: { ...state, ...body } })
    }

    const event: TLambdaEvent = {
      queryStringParameters: { referenceId },
      context: {
        signature: 'encrypted',
        authAccess: {
          bearerBaseURL: 'https://void.bearer.sh',
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

describe('RetrieveState', () => {
  describe('existing data', () => {
    const referenceId = 'SPONGE_BOB'

    const event: TLambdaEvent = {
      queryStringParameters: { referenceId },
      context: {
        bearerBaseURL: 'https://void.bearer.sh',
        signature: 'encrypted',
        authAccess: {
          bearerBaseURL: 'wedontcare'
        }
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
      mockAxios.mockResponse({ data: { Item: { title: 'Sponge bob', referenceId } } })
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
