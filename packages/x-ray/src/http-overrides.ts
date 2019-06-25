// Mostly inspired from: https://github.com/meetearnest/global-request-logger/blob/master/index.js

import urlParser from 'url'
import http from 'http'
import https from 'https'

import log from './logger'
import { _HANDLER, STAGE } from './constants'

const EXCLUDED_API = new Set(['logs.eu-west-3.amazonaws.com', 'logs.eu-west-1.amazonaws.com'])

const buildBillingPayload = (info: TRequestPayload) => {
  const { host, pathname, method } = info.request

  return {
    message: {
      pathname,
      method,
      path: host,
      intentName: _HANDLER,
      clientId: process.env.clientId,
      integrationUuid: process.env.scenarioUuid,
      stage: STAGE,
      type: 'externalCall',
      tid: process.env._X_AMZN_TRACE_ID
    },
    timestamp: Date.now()
  }
}

export const overrideRequestMethod = (
  module: typeof http | typeof https,
  billingLogger: (...arg: any[]) => void = log.extend('externalCall'),
  userLogger: (...arg: any[]) => void = log.extend('user')
) => {
  // override http.request method
  const _request = module.request

  function request(...args: any): http.ClientRequest {
    // @ts-ignore
    const req = _request(...arguments)
    const info: TRequestPayload = {
      request: {
        method: req.method || 'get',
        headers: { ...req.getHeaders() }
      } as any,
      response: {} as any
    }
    const url: string | urlParser.UrlWithParsedQuery = args[0]

    let parsedUrl: urlParser.UrlWithParsedQuery

    if (typeof url === 'string') {
      parsedUrl = urlParser.parse(url, true)
    } else {
      parsedUrl = url
    }
    const { port, path, host, protocol, auth, hostname, hash, search, query, pathname, href } = parsedUrl
    info.request = {
      ...info.request,
      port,
      path,
      host,
      protocol,
      auth,
      hostname,
      hash,
      search,
      query,
      pathname,
      href,
      url: `${protocol}//${host || hostname}${path}`
    }
    const requestData: string[] = []
    const requestStart = Date.now()
    const originalWrite = req.write

    req.write = function() {
      logBodyChunk(requestData, arguments[0])
      originalWrite.apply(req, arguments)
    }

    req.on('error', (error: any) => {
      info.request.error = error
    })

    req.on('response', (res: http.ServerResponse) => {
      info.request.body = requestData.join('')
      info.response = {
        status: res.statusCode,
        headers: (res as any).headers
      }

      const responseData: string[] = []
      res.on('data', (data: string) => {
        logBodyChunk(responseData, data)
      })

      res.on('end', () => {
        if (!EXCLUDED_API.has(info.request.hostname!)) {
          info.response.body = responseData.join('')
          info.duration = Date.now() - requestStart

          billingLogger('%j', buildBillingPayload(info))
          userLogger('%j', {
            buid: process.env.scenarioUuid,
            payloadType: 'EXTERNAL_CALL',
            payload: info,
            service: 'function',
            tid: process.env._X_AMZN_TRACE_ID
          })
        }
      })
      res.on('error', (error: any) => {
        info.response.error = error
      })
    })

    return req
  }

  function get(...args: any): http.ClientRequest {
    // @ts-ignore
    const req = request(...arguments)
    req.end()
    return req
  }

  module.request = request
  module.get = get
}

type TRequestPayload = {
  request: { url: string; method: string; headers: any; error?: any; body?: any } & urlParser.UrlWithParsedQuery
  response: {
    body?: any
    status: number
    headers: any
    error?: any
  }
  duration?: number
}
const maxBodyLength = 2024 * 1000 * 3

function logBodyChunk(array: any[], chunk: string) {
  if (chunk) {
    let toAdd = chunk
    const newLength = array.length + chunk.length
    if (newLength > maxBodyLength) {
      toAdd = chunk.slice(0, maxBodyLength - newLength)
    }
    array.push(toAdd)
  }
}
