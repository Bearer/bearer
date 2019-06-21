// Mostly inspired from: https://github.com/meetearnest/global-request-logger/blob/master/index.js

import url from 'url'
import http from 'http'
import https from 'https'

import logger from './logger'
import { _HANDLER } from './constants'

const EXCLUDED_API = new Set(['logs.eu-west-3.amazonaws.com', 'logs.eu-west-1.amazonaws.com'])

export const overrideRequestMethod = (module: typeof http | typeof https) => {
  // override http.request method
  const _request = module.request

  function request(...args: any): http.ClientRequest {
    // @ts-ignore
    const req = _request(...arguments)
    const info = {
      request: {
        method: req.method || 'get',
        headers: { ...req.getHeaders() }
      } as { method: string; headers: any; error?: any; body: any } & url.UrlWithStringQuery,
      response: {} as {
        body: any
        statusCode: number
        headers: any
        error?: any
      }
    }
    const requestUrl = args[0]

    if (requestUrl === 'string') {
      const { port, path, host, protocol, auth, hostname, hash, search, query, pathname, href } = url.parse(requestUrl)
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
        href
      }
    }
    const requestData: string[] = []
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
      info.response.statusCode = res.statusCode
      info.response.headers = (res as any).headers

      const responseData: string[] = []
      res.on('data', (data: string) => {
        logBodyChunk(responseData, data)
      })

      res.on('end', () => {
        if (!EXCLUDED_API.has(info.request.host!)) {
          info.response.body = responseData.join('')
          logger.extend('externalCall')('%j', {
            buid: process.env.scenarioUuid,
            payloadType: 'FUNCTION',
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
