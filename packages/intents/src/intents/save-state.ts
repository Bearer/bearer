import debug from '@bearer/logger'
import uuid from 'uuid/v1'

import * as d from '../declaration'
import { DBClient as CLIENT } from '../db-client'
import { bodyFromEvent, eventAsActionParams } from './utils'

const logger = debug('intents:fetch-state')

export class SaveState {
  static intentPromise(action: d.ISaveStateAction) {
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

  // TODO: remove as soon as async intents are released
  static intent(action: d.ISaveStateIntentAction) {
    // tslint:disable-next-line:variable-name
    const DBClient = CLIENT.instance

    return (event: d.TLambdaEvent, _context: any, lambdaCallback: d.TLambdaCallback): void => {
      const referenceId = event.queryStringParameters.referenceId
      const dbClient = DBClient(event.context.signature)
      try {
        dbClient
          .getData(referenceId)
          .then(savedState => {
            const state = savedState ? savedState.Item : {}
            try {
              action(
                event.context,
                event.queryStringParameters,
                bodyFromEvent(event),
                state,
                (result: { state: any; data?: any }) => {
                  if (savedState || referenceId) {
                    dbClient
                      .updateData(referenceId, result.state)
                      .then(() => {
                        lambdaCallback(null, { meta: { referenceId }, data: result.data })
                      })
                      .catch(error => lambdaCallback(error.toString(), { error: error.toString() }))
                  } else {
                    dbClient
                      .saveData(result.state)
                      .then(data => {
                        lambdaCallback(null, { meta: { referenceId: data.Item.referenceId }, data: result.data })
                      })
                      .catch(error => lambdaCallback(error.toString(), { error: error.toString() }))
                  }
                }
              )
            } catch (error) {
              lambdaCallback(error.toString(), { error: error.toString() })
            }
          })
          .catch(error => lambdaCallback(error.toString(), { error: error.toString() }))
      } catch (error) {
        lambdaCallback(error.toString(), { error: error.toString() })
      }
    }
  }
}

export class SaveStateActionExecutionError extends Error {}
export class SaveStateSavingStateError extends Error {}
