import { SaveState, RetrieveState } from '../src/index'
import mockAxios from 'jest-mock-axios'

afterEach(() => {
  mockAxios.reset()
})

describe('index', () => {
  it('export SaveState', () => {
    expect(SaveState).toBeTruthy()
  })

  it('export RetrieveState', () => {
    expect(RetrieveState).toBeTruthy()
  })
})

describe('SaveState', () => {
  describe('.intent', () => {
    const referenceId = 'SPONGE_BOB'

    const action = (_context: any, _params: any, body: any, state: any, callback: (state: any) => void) => {
      callback({ ...state, ...body })
    }

    const event = {
      queryStringParameters: { referenceId },
      context: { bearerBaseURL: 'http://void.bearer.sh' },
      body: {
        pullRequests: []
      }
    }

    describe('reference exists', () => {
      const savedData = {
        data: 'server says hello!'
      }

      function mockRetrieveSucces() {
        expect(mockAxios.get).toHaveBeenCalledWith('api/v1/items/SPONGE_BOB')
        const responseObj = { data: { Item: savedData } }
        mockAxios.mockResponse(responseObj)
      }

      it('retrieves reference, update it and return payload', () => {
        const lambdaCallback = jest.fn()

        SaveState.intent(action)(event, {}, lambdaCallback)

        mockRetrieveSucces()

        expect(mockAxios.put).toHaveBeenCalledWith('api/v1/items/SPONGE_BOB', {
          ReadAllowed: true,
          ...savedData,
          pullRequests: []
        })
        mockAxios.mockResponse({}) // data is not used

        expect(lambdaCallback).toHaveBeenCalledWith(null, {
          meta: { referenceId },
          data: { ...savedData, pullRequests: [] }
        })
      })

      it('retrieves reference, fails update and return error', () => {
        const lambdaCallback = jest.fn()

        SaveState.intent(action)(event, {}, lambdaCallback)

        mockRetrieveSucces()

        expect(mockAxios.put).toHaveBeenCalledWith('api/v1/items/SPONGE_BOB', {
          ReadAllowed: true,
          ...savedData,
          pullRequests: []
        })
        mockAxios.mockError({ status: 500 })

        expect(lambdaCallback).toHaveBeenCalledWith(expect.any(String), { error: 'Cannot update data' })
      })

      it('retrieves data, action throw error, returns error', () => {
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
        expect(mockAxios.get).toHaveBeenCalledWith('api/v1/items/SPONGE_BOB')
        const responseObj = { data: { Item: null }, status: 404 }
        mockAxios.mockError(responseObj)
      }

      it('create reference and return payload', () => {
        const lambdaCallback = jest.fn()

        SaveState.intent(action)(event, {}, lambdaCallback)

        mockRetrieveFailure()

        expect(mockAxios.post).toHaveBeenCalledWith('api/v1/items', { ReadAllowed: true, pullRequests: [] })
        const responseObj = { data: { Item: { referenceId: 'PATRICK' } } }
        mockAxios.mockResponse(responseObj)

        expect(mockAxios.post).toHaveBeenCalledWith()
        expect(lambdaCallback).toHaveBeenCalledWith(null, { meta: { referenceId: 'PATRICK', data: { ...event.body } } })
      })

      it('fails create reference and return error', () => {
        const lambdaCallback = jest.fn()

        SaveState.intent(action)(event, {}, lambdaCallback)

        mockRetrieveFailure()

        expect(mockAxios.post).toHaveBeenCalledWith('api/v1/items', { ReadAllowed: true, pullRequests: [] })
        const responseObj = {}
        mockAxios.mockError(responseObj)
        expect(lambdaCallback).toHaveBeenCalledWith(expect.any(String), { error: 'Cannot save data' })
      })
      it('action throw error, returns error', () => {
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
      it('returns error', () => {
        const lambdaCallback = jest.fn()

        SaveState.intent(action)(event, {}, lambdaCallback)
        expect(mockAxios.get).toHaveBeenCalledWith('api/v1/items/SPONGE_BOB')
        const responseObj = { data: { Item: null }, status: 500 }
        mockAxios.mockError(responseObj)

        expect(mockAxios.post).not.toHaveBeenCalledWith()
        expect(mockAxios.put).not.toHaveBeenCalledWith()
        expect(lambdaCallback).toHaveBeenCalledWith(null, { error: 'Error whild trying to fetch current state' })
      })
    })
  })
})
