import debug from '@bearer/logger'
import { bearerOverride, setupFunctionIdentifiers } from '@bearer/x-ray'

import * as d from '../declaration'
import { eventAsActionParams, BACKEND_ONLY_ERROR } from './utils'
import { Store } from '../store'

const logger = debug('functions:fetch-state')

interface FetchDataImplementation<T extends FetchData> {
  new (): T
}

export abstract class FetchData<ReturnedData = any, TError = any, AuthContext = d.TAuthContext> {
  static serverSideRestricted = false

  // expected implementation
  abstract async action(
    event: d.TFetchActionEvent<any, AuthContext, any>
  ): Promise<d.TFetchPayload<ReturnedData, TError>>

  // Internal
  static call(aPrototype: FetchDataImplementation<any>) {
    const action = new aPrototype.prototype.constructor().action as d.TFetchAction
    const requiresBackend = (aPrototype as any).serverSideRestricted

    return async (event: d.TLambdaEvent) => {
      if (requiresBackend && !event.context.isBackend) {
        return { error: BACKEND_ONLY_ERROR }
      }

      const updatedEvent = Object.assign({}, event)
      updatedEvent.context.integrationUuid = event.context.buid
      setupFunctionIdentifiers(updatedEvent)

      const functionEvent: d.TFetchActionEvent = {
        ...eventAsActionParams(event),
        store: new Store(event.context.signature)
      }

      try {
        const { statusCode, error, data, referenceId }: d.TFetchPayload<any, any> = await action(functionEvent)
        if (error) {
          logger(error)
          return { statusCode, error }
        }
        return { statusCode, data, referenceId }
      } catch (error) {
        logger.extend('FetchActionExecutionError')(error)
        throw new FetchActionExecutionError(error)
      }
    }
  }

  static init() {
    bearerOverride()
    return FetchData.call(this as any)
  }
}

export class FetchActionExecutionError extends Error {}
