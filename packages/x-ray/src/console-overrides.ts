import { STAGE, _HANDLER } from './constants'
import uuid = require('uuid')
import console from 'console'

import log from './logger'

export const overrideConsole = (module: typeof console, userLogger: any = log.extend('user')) => {
  const _log = module.log

  module.log = (args: any) => {
    userLogger('%j', buildLog(args))
    return _log.apply(module, args)
  }
}

const buildLog = (args: any) => {
  return {
    buid: process.env.scenarioUuid,
    payloadType: 'DEBUG',
    service: 'function',
    payload: {
      log: args.toString(),
      intentName: _HANDLER,
      uuid: uuid()
    },
    stage: STAGE,
    tid: process.env._X_AMZN_TRACE_ID
  }
}
