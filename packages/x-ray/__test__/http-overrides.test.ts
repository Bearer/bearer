import { overrideGetMethod, overrideRequestMethod } from '../src/http-overrides'
import http from 'http'
import { sendToCloudwatchGroup } from '../src/cloud-watch-logs'
import { httpClient } from './helpers/utils'

jest.mock('../src/constants')
jest.mock('../src/cloud-watch-logs')
process.env.clientId = '132464737464748494404949984847474848'
process.env.scenarioUuid = 'scenarioUuid'

describe('override http', () => {
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
        data.emit('end')
      })
    })

    expect(sendToCloudwatchGroup).toHaveBeenCalledWith({
      message: {
        clientId: '132464737464748494404949984847474848',
        integrationUuid: 'scenarioUuid',
        intentName: 'MyHandler',
        method: 'GET',
        path: 'www.google.com',
        pathname: '/',
        responseStatus: 200,
        responseStatusMesage: 'OK',
        stage: 'test'
      },
      timestamp: expect.any(Number)
    })
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
      http.get('http://www.google.com/', (data: http.IncomingMessage) => {
        res(data)
        data.emit('end')
      })
    })

    expect(sendToCloudwatchGroup).toHaveBeenCalledWith({
      message: {
        clientId: '132464737464748494404949984847474848',
        integrationUuid: 'scenarioUuid',
        intentName: 'MyHandler',
        method: 'GET',
        path: 'www.google.com',
        pathname: '/',
        responseStatus: 200,
        responseStatusMesage: 'OK',
        stage: 'test'
      },
      timestamp: expect.any(Number)
    })
  })
})
