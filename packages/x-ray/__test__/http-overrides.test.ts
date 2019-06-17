import http from 'http'
import { overrideGetMethod, overrideRequestMethod } from '../src/http-overrides'
import { httpClient } from './helpers/utils'
import logger from '../src/logger'
import uuid from 'uuid'

jest.mock('uuid')
jest.mock('../src/constants')
process.env.clientId = '132464737464748494404949984847474848'
process.env.scenarioUuid = 'scenarioUuid'
process.env._X_AMZN_TRACE_ID = 'Root=is-uuid;Parent=my-uuid'

const spy = jest.spyOn(logger, 'extend')
const payloadUUID = 'test-uuid'

// @ts-ignore
uuid.mockImplementation(() => payloadUUID)

beforeEach(() => {
  spy.mockClear()
})

describe('override http', () => {
  const log = jest.fn()

  beforeEach(() => {
    log.mockClear()
    spy.mockImplementationOnce(() => log as any)
  })

  it('overrides request method', () => {
    // @ts-ignore
    expect(httpClient._request).toBeUndefined()
    // @ts-ignore
    overrideRequestMethod(httpClient)
    // @ts-ignore
    expect(httpClient._request).toBeDefined()
  })

  it('traces request call with the right payload', async () => {
    // @ts-ignore
    expect(httpClient._request).toBeDefined()

    await new Promise((res, _rej) => {
      httpClient.request('http://www.google.com/', (data: http.IncomingMessage) => {
        res(data)
        data.resume()
      })
    })

    expect(log.mock.calls[0][1]).toMatchSnapshot({ timestamp: expect.any(Number) })
  })

  it('overrides the get method', () => {
    // @ts-ignore
    expect(httpClient._get).toBeUndefined()
    overrideGetMethod(httpClient)
    // @ts-ignore
    expect(httpClient._get.name).toEqual('get')
    // @ts-ignore
    expect(httpClient._get).toBeDefined()
  })

  it('traces get call with the right payload', async () => {
    // @ts-ignore
    expect(httpClient.get).toBeDefined()
    await new Promise((res, _rej) => {
      httpClient.get('http://www.google.com/', (data: http.IncomingMessage) => {
        res(data)
        data.resume()
      })
    })

    expect(log.mock.calls[0][1]).toMatchSnapshot({ timestamp: expect.any(Number) })
  })
})

describe('log filtering', () => {
  it('is calling logger for non aws calls', async () => {
    await new Promise((res, _rej) => {
      httpClient.request('http://www.google.com/', (data: http.IncomingMessage) => {
        res(data)
        data.resume()
      })
    })

    expect(spy).toHaveBeenCalledWith('externalCall')
  })

  it('is not calling logger for eu-west-3 aws logs calls', async () => {
    await new Promise((res, _rej) => {
      httpClient.request('http://logs.eu-west-3.amazonaws.com/', (data: http.IncomingMessage) => {
        res(data)
        data.resume()
      })
    })

    expect(spy).not.toHaveBeenCalled()
  })

  it('is not calling logger for eu-west-1 aws logs calls', async () => {
    await new Promise((res, _rej) => {
      httpClient.request('http://logs.eu-west-1.amazonaws.com/', (data: http.IncomingMessage) => {
        res(data)
        data.resume()
      })
    })

    expect(spy).not.toHaveBeenCalled()
  })
})
