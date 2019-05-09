import { captureHttps, setupFunctionIdentifiers } from '../src/index'
import { httpClient, event } from './helpers/utils'
import http from 'http'

jest.mock('../src/constants')

describe('captureHttp', () => {
  it('captures https without identifiers', async () => {
    captureHttps(httpClient)
    expect((httpClient as any)._request).toBeDefined()
    await new Promise((res, _rej) => {
      httpClient.request('http://www.google.com/', (data: http.IncomingMessage) => {
        res(data)
        data.resume()
      })
    })

    expect(process.env.clientId).toEqual(undefined)
    expect(process.env.scenarioUuid).toEqual(undefined)
  })

  it('captures https with identifiers', async () => {
    const functionEvent = event()
    captureHttps(httpClient)
    setupFunctionIdentifiers(functionEvent)
    await new Promise((res, _rej) => {
      httpClient.request('http://www.google.com/', (data: http.IncomingMessage) => {
        res(data)
        data.resume()
      })
    })

    expect(process.env.clientId).toEqual(functionEvent.context.clientId)
    expect(process.env.scenarioUuid).toEqual(functionEvent.context.integrationUuid)
  })
})

describe('setupFunctionIdentifiers', () => {
  it('setup all function identifiers', () => {
    const functionEvent = event()
    setupFunctionIdentifiers(functionEvent)
    expect(process.env.clientId).toEqual(functionEvent.context.clientId)
    expect(process.env.scenarioUuid).toEqual(functionEvent.context.integrationUuid)
  })
})
