import debug from '@bearer/logger'

import * as d from '../declaration'
import { eventAsActionParams } from './utils'

const logger = debug('intents:fetch-state')

export abstract class FetchData<ReturnedData = any, Error = any, AuthContext = any> {
  // expected implementation
  abstract async action<Params = any>(
    event: d.TFetchActionEvent<AuthContext, Params>
  ): Promise<d.TFetchPayload<ReturnedData, Error>>

  // Internal
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

  static init() {
    return FetchData.intent(new (this.prototype.constructor as any)().action)
  }
}

export class FetchActionExecutionError extends Error {}
