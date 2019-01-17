import debug from '@bearer/logger'

import * as d from '../declaration'
import { bodyFromEvent, fetchData, eventAsActionParams } from './utils'

const logger = debug('intents:fetch-state')

export class FetchData {
  // TODO: remove as soon as async intents are released
  static intent(action: d.TFetchDataAction) {
    return (event: d.TLambdaEvent, _context, lambdaCallback: d.TLambdaCallback) => {
      action(event.context, event.queryStringParameters, bodyFromEvent(event), result => {
        fetchData(lambdaCallback, result)
      })
    }
  }

  static intentPromise(action: d.TFetchAction) {
    return async (event: d.TLambdaEvent) => {
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
}

export class FetchActionExecutionError extends Error {}
