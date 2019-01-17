import debug from '@bearer/logger'

import * as d from '../declaration'
import { eventAsActionParams } from './utils'

const logger = debug('intents:fetch-state')

export class FetchData {
  static intent(action: d.TFetchAction) {
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
