import debug from '@bearer/logger'
const logger = debug('intents')
import { sendToCloudwatchGroup } from './cloud-watch-logs'
import { STAGE, _HANDLER } from './constants'

export const overrideRequestMethod = (module: any) => {
  // override http.request method
  module._request = module.request
  module.request = (options: any, callaback: any) => {
    if (options.host === 'logs.eu-west-3.amazonaws.com') {
      console.log('ignore logs.eu-west-3.amazonaws.com requests')
      return module._request(options, callaback)
    }

    console.log(process.env)
    const payload: any = {
      message: {
        path: options.hostname || options.host,
        pathname: options.pathname || options.path,
        method: options.method,
        intentName: _HANDLER,
        clientId: process.env.clientId,
        integrationUuid: process.env.scenarioUuid,
        stage: STAGE
      },
      timestamp: new Date().getTime()
    }

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
    const payload: any = {
      message: {
        path: options.hostname || options.host,
        pathname: options.pathname || options.path,
        method: options.method,
        intentName: _HANDLER,
        clientId: process.env.clientId,
        integrationUuid: process.env.scenarioUuid,
        stage: STAGE
      },
      timestamp: new Date().getTime()
    }

    return module._get(options, (res: any) => {
      res.on('end', async () => {
        payload.responseStatus = res.statusCode
        payload.responseStatusMesage = res.statusMessage

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
