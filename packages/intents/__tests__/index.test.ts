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
    it('does something', () => {
      SaveState.intent(() => {})({ queryStringParameters: {}, context: {} }, {}, jest.fn())
      console.log('[BEARER]', 'axios.lastReqGet()', mockAxios.lastReqGet())
      // expect(mockAxios.get).toHaveBeenCalledWith('/.*/', { data: {} })
      // simulating a server response
      // let responseObj = { data: 'server says hello!' }
      // mockAxios.mockResponse(responseObj)

      expect(true).toBeTruthy()
    })
  })
})
