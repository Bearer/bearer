import debug from '@bearer/logger'
import uuid from 'uuid/v1'

import * as d from '../declaration'
import { DBClient as CLIENT } from '../db-client'
import { eventAsActionParams } from './utils'

const logger = debug('intents:fetch-state')

export abstract class SaveState<State = any, ReturnedData = any, Error = any, AuthContext = any> {
  // expected implementation
  abstract async action(
    event: d.TSaveActionEvent<State, any, AuthContext>
  ): Promise<d.TSaveStatePayload<State, ReturnedData, Error>>

  // Internal
  static intent(action: d.TSaveStateAction) {
    // tslint:disable-next-line:variable-name
    const DBClient = CLIENT.instance

    return async (event: d.TLambdaEvent) => {
      const providedReferenceId = event.queryStringParameters.referenceId
      const dbClient = DBClient(event.context.signature)

      try {
        const savedData = await dbClient.getData(providedReferenceId)
        const currentState = savedData ? savedData.Item : {}
        const { state, data, error } = await action({ ...eventAsActionParams(event), state: currentState })
        if (error) {
          return { error }
        }
        const referenceId = providedReferenceId || uuid()
        try {
          await dbClient.updateData(referenceId, state)
        } catch (error) {
          logger.extend('SaveStateSavingStateError')(error)
          throw new SaveStateSavingStateError()
        }
        return { data, meta: { referenceId } }
      } catch (error) {
        logger.extend('SaveStateActionExecutionError')(error)
        throw new SaveStateActionExecutionError()
      }
    }
  }

  static init() {
    return SaveState.intent(new (this.prototype.constructor as any)().action)
  }
}

export class SaveStateActionExecutionError extends Error {}
export class SaveStateSavingStateError extends Error {}
