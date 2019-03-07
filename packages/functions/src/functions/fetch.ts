import debug from '@bearer/logger'

import * as d from '../declaration'
import { eventAsActionParams } from './utils'

const logger = debug('functions:fetch-state')

export abstract class FetchData<ReturnedData = any, TError = any, AuthContext = any> {
  // expected implementation
  abstract async action(
    event: d.TFetchActionEvent<any, AuthContext, any>
  ): Promise<d.TFetchPayload<ReturnedData, TError>>

  // Internal
  static call(action: d.TFetchAction) {
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
    return FetchData.call(new (this.prototype.constructor as any)().action)
  }
}

export class FetchActionExecutionError extends Error {}
