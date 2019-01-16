import * as d from '../declaration'
import { DBClient as CLIENT } from '../db-client'
import { bodyFromEvent } from './utils'

export class SaveState {
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
