import debug from '@bearer/logger'

import * as d from '../declaration'
import { bodyFromEvent, fetchData, paramsFromEvent } from './utils'

const logger = debug('intents:fetch-state')

export class FetchData {
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
        const { error, data }: d.TFetchPayload<any, any> = await action({
          context: event.context,
          params: paramsFromEvent(event)
        })

        if (error) {
          logger(error)
          return { error }
        }
        return { data }
      } catch (error) {
        logger.extend('ActionExecutionError')(error)
        throw new ActionExecutionError(error)
      }
    }
  }
}

export class ActionExecutionError extends Error {}
