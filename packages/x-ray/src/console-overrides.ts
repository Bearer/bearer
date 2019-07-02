import { STAGE, _HANDLER } from './constants'
import uuid = require('uuid')
import console from 'console'

import log from './logger'
import stream from 'stream'

const out = new stream.Writable()

out._write = (chunk, _encoding, next) => {
  userLogger('%j', buildLog(chunk.toString()))
  next()
}

const userLogger = log.extend('user')

const myConsole = new console.Console(out)

export const overrideConsole = (module: typeof console) => {
  const _log = module.log

  module.log = (...args: any) => {
    myConsole.log(args)

    return _log.apply(module, args)
  }
}

const buildLog = (args: any) => {
  return {
    buid: process.env.scenarioUuid,
    payloadType: 'DEBUG',
    service: 'function',
    payload: {
      log: args,
      intentName: _HANDLER,
      uuid: uuid()
    },
    stage: STAGE,
    tid: process.env._X_AMZN_TRACE_ID
  }
}
