import logger from './logger'
import { sendToCloudwatchGroup } from './cloud-watch-logs'
import { STAGE, _HANDLER } from './constants'
import url from 'url'

export const overrideRequestMethod = (module: any) => {
  // override http.request method
  module._request = module.request
  module.request = (options: any, callaback: any) => {
    const settings = parseOptions(options)

    if (settings.host === 'logs.eu-west-3.amazonaws.com') {
      logger('ignore logs.eu-west-3.amazonaws.com requests')
      return module._request(options, callaback)
    }

    const payload: any = buildTracePatload(settings)

    return module._request(options, (res: any) => {
      res.on('end', async () => {
        payload.message.responseStatus = res.statusCode
        payload.message.responseStatusMesage = res.statusMessage

        if (!STAGE) {
          logger('%j', payload)
        } else {
          await sendToCloudwatchGroup(payload)
        }
      })

      if (typeof callaback === 'function') {
        callaback(res)
      } else {
        res.resume()
      }
    })
  }
}

export const overrideGetMethod = (module: any) => {
  // overrride the http.get method
  module._get = module.get
  module.get = (options: any, callaback: any) => {
    const settings = parseOptions(options)

    if (settings.host === 'logs.eu-west-3.amazonaws.com') {
      logger('ignore logs.eu-west-3.amazonaws.com requests')
      return module._get(options, callaback)
    }

    const payload: any = buildTracePatload(settings)
    return module._get(options, (res: any) => {
      res.on('end', async () => {
        payload.message.responseStatus = res.statusCode
        payload.message.responseStatusMesage = res.statusMessage

        if (!STAGE) {
          logger('%j', payload)
        } else {
          await sendToCloudwatchGroup(payload)
        }
      })
      if (typeof callaback === 'function') {
        callaback(res)
      } else {
        res.resume()
      }
    })
  }
}

const parseOptions = (options: any) => {
  if (typeof options === 'string') {
    return url.parse(options)
  }

  return options
}

const buildTracePatload = (settings: any) => {
  return {
    message: {
      path: settings.hostname || settings.host,
      pathname: settings.pathname || settings.path,
      method: settings.method || 'GET',
      intentName: _HANDLER,
      clientId: process.env.clientId,
      integrationUuid: process.env.scenarioUuid,
      stage: STAGE
    },
    timestamp: new Date().getTime()
  }
}
