import logger from './logger'
import { STAGE, _HANDLER, _X_AMZN_TRACE_ID } from './constants'
import url from 'url'
import uuid = require('uuid')

const EXCLUDED_API = ['logs.eu-west-3.amazonaws.com', 'logs.eu-west-1.amazonaws.com']

export const overrideRequestMethod = (module: any) => {
  // override http.request method
  module._request = module.request
  module.request = (options: any, callback: any) => {
    const settings = parseOptions(options)
    const payload: any = buildTracePayload(settings)
    return baseRequest(module, options, payload, settings, callback)
  }
}

export const overrideGetMethod = (module: any) => {
  // override the http.get method
  module._get = module.get
  module.get = (options: any, callback: any) => {
    const settings = parseOptions(options)
    const payload: any = buildTracePayload(settings)
    return baseGet(module, options, payload, settings, callback)
  }
}

const parseOptions = (options: any) => {
  if (typeof options === 'string') {
    return url.parse(options)
  }

  return options
}

const buildTracePayload = (settings: any) => {
  return {
    message: {
      path: settings.hostname || settings.host,
      pathname: settings.pathname || settings.path,
      method: settings.method || 'GET',
      intentName: _HANDLER,
      clientId: process.env.clientId,
      integrationUuid: process.env.scenarioUuid,
      stage: STAGE,
      uuid: uuid(),
      type: 'externalCall',
      tid: _X_AMZN_TRACE_ID
    },
    timestamp: new Date().getTime()
  }
}
function baseGet(module: any, options: any, payload: any, settings: any, callback: any) {
  const req = module
    ._get(options, async (res: any) => {
      onRequestEnd(res, payload, settings, callback)
    })
    .on('err', async (e: { message: string }) => {
      logger('error event received')
      logger('response error %s', e.message)
    })
  return req
}

function baseRequest(module: any, options: any, payload: any, settings: any, callback: any) {
  const req = module
    ._request(options, async (res: any) => {
      onRequestEnd(res, payload, settings, callback)
    })
    .on('err', async (e: { message: string }) => {
      logger('error event received')
      logger('response error %s', e.message)
    })
  return req
}

const onRequestEnd = (res: any, payload: any, settings: any, callback: any) => {
  res.on('end', () => {
    payload.message.responseStatus = res.statusCode
    payload.message.responseStatusMessage = res.statusMessage
    if (!EXCLUDED_API.includes(settings.host)) {
      logger.extend('externalCall')('%j', payload)
    }
  })

  if (typeof callback === 'function') {
    callback(res)
  } else if (res && res.listenerCount('end') === 1) {
    res.resume()
  }
}
