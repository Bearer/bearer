import EventEmitter from 'events'

export class TestEmitter extends EventEmitter {
  method: string = 'GET'
  url: string = '/'
  connection: any = {}
  req: any = {}
  statusCode: number = 200
  statusMessage: string = 'KO'

  resume() {
    this.emit('resume')
  }
}

function buildFakeRequest() {
  const request = new TestEmitter()
  request.method = 'GET'
  request.url = '/'
  request.connection = { remoteAddress: 'myhost' }
  return request
}

function buildFakeResponse() {
  const response = new TestEmitter()
  response.statusCode = 200
  response.statusMessage = 'OK'

  response.resume = () => {
    response.emit('end')
  }
  return response
}

export const httpClient = {
  request: (options: any, callback: Function) => {
    const fakeRequest = buildFakeRequest()
    const fakeReponse = buildFakeResponse()
    fakeReponse.req = fakeRequest
    callback(fakeReponse)
    return fakeRequest
  },
  get: (options: any, callback: Function) => {
    const fakeRequest = buildFakeRequest()
    const fakeReponse = buildFakeResponse()
    fakeReponse.req = fakeRequest
    callback(fakeReponse)
    return fakeRequest
  }
}

export const expectedResponse = {
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
}
