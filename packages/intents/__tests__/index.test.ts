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

    it('retrieves state and update it', () => {
      const savedData = {
        data: 'server says hello!'
      }
      const lambdaCallback = jest.fn()
      const action = (_context: any, _params: any, body: any, state: any, callback: (state: any) => void) => {
        callback({ ...state })
      }

      SaveState.intent(action)(
        { queryStringParameters: { referenceId }, context: { bearerBaseURL: 'http://void.bearer.sh' } },
        {},
        lambdaCallback
      )
      expect(mockAxios.get).toHaveBeenCalledWith('api/v1/items/SPONGE_BOB')
      const responseObj = { data: { Item: savedData } }
      mockAxios.mockResponse(responseObj)

      expect(mockAxios.put).toHaveBeenCalledWith('api/v1/items/SPONGE_BOB', { ReadAllowed: true, ...savedData })
      mockAxios.mockResponse({})

      expect(lambdaCallback).toHaveBeenCalledWith(null, { meta: { referenceId }, data: savedData })
    })
  })
})
