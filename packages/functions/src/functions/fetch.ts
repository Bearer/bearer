import debug from '@bearer/logger'
import http from 'http'
import https from 'https'
import { captureHttps } from '@bearer/x-ray'

import * as d from '../declaration'
import { eventAsActionParams } from './utils'

const logger = debug('functions:fetch-state')

interface FetchDataImplementation<T extends FetchData> {
  new (): T
}

export abstract class FetchData<ReturnedData = any, TError = any, AuthContext = any> {
  static backendOnly = false

  // expected implementation
  abstract async action(
    event: d.TFetchActionEvent<any, AuthContext, any>
  ): Promise<d.TFetchPayload<ReturnedData, TError>>

  // Internal
  static call(aPrototype: FetchDataImplementation<any>) {
    const action = new aPrototype.prototype.constructor().action as d.TFetchAction
    const requiresBackend = (aPrototype as any).backendOnly

    return async (event: d.TLambdaEvent) => {
      if (requiresBackend && !event.context.isBackend) {
        return {
          error: {
            code: 'UNAUTHORIZED_FUNCTION_CALL',
            message: "This function can't be called"
          }
        }
      }
      captureHttps(http, event)
      captureHttps(https, event)
      try {
        const { error, data }: d.TFetchPayload<any, any> = await action(eventAsActionParams(event))
        if (error) {
          logger(error)
          return { error }
        }
        return { data }
      } catch (error) {
        logger.extend('FetchActionExecutionError')(error)
        throw new FetchActionExecutionError(error)
      }
    }
  }

  static init() {
    return FetchData.call(this as any)
  }
}

export class FetchActionExecutionError extends Error {}
