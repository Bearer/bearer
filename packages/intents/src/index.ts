import * as d from './declaration'
export * from './declaration'

import { DBClient as CLIENT } from './DBClient'
import { sendErrorMessage, sendSuccessMessage } from './lambda'

export const DBClient = CLIENT.instance

// tslint:disable-next-line:no-unnecessary-class
export class Intent {
  static fetchData(callback: d.TLambdaCallback, { data, error }: { data?: any; error?: any }) {
    if (data) {
      sendSuccessMessage(callback, { data })
    } else {
      sendErrorMessage(callback, { error: error || 'Unkown error' })
    }
  }
}

// tslint:disable-next-line:no-unnecessary-class
class BaseIntent {
  static get display(): string {
    throw new Error('Extending class needs to implement `static intent(action)` method')
  }

  static intent(_action): (event: any, context: any, callback: (...args: any[]) => any) => any {
    throw new Error('Extending class needs to implement `static intent(action)` method')
  }
}

class GenericIntentBase extends BaseIntent {
  static get isStateIntent(): boolean {
    return false
  }

  static get isGlobalIntent(): boolean {
    return true
  }
}

class StateIntentBase extends BaseIntent {
  static get isStateIntent(): boolean {
    return true
  }

  static get isGlobalIntent(): boolean {
    return false
  }
}

export class SaveState extends StateIntentBase {
  static get display(): string {
    return 'SaveState'
  }

  static intent(action: d.ISaveStateIntentAction) {
    return (event: d.TLambdaEvent, _context: any, lambdaCallback: d.TLambdaCallback): void => {
      const { referenceId } = event.queryStringParameters
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
                (result: { state: any; data: any }) => {
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

export class RetrieveState extends StateIntentBase {
  static get display(): string {
    return 'RetrieveState'
  }

  static intent(action: d.IRetrieveStateIntentAction) {
    return (event: d.TLambdaEvent, _context: any, lambdaCallback: d.TLambdaCallback): void => {
      const { referenceId } = event.queryStringParameters
      try {
        DBClient(event.context.signature)
          .getData(referenceId)
          .then(state => {
            if (state) {
              action(
                event.context,
                event.queryStringParameters,
                state.Item,
                (recoveredState: { data: any }): void => {
                  lambdaCallback(null, { meta: { referenceId: state.Item.referenceId }, data: recoveredState.data })
                }
              )
            } else {
              lambdaCallback(null, { statusCode: 404, body: JSON.stringify({ error: 'No data found', referenceId }) })
            }
          })
      } catch (error) {
        lambdaCallback(error.toString(), { error: error.toString() })
      }
    }
  }
}

export class FetchData extends GenericIntentBase {
  static get display(): string {
    return 'FetchData'
  }

  static intent(action: d.TFetchDataAction) {
    return (event: d.TLambdaEvent, _context, lambdaCallback: d.TLambdaCallback) => {
      action(event.context, event.queryStringParameters, bodyFromEvent(event), result => {
        Intent.fetchData(lambdaCallback, result)
      })
    }
  }
}

function bodyFromEvent(event: d.TLambdaEvent): any {
  const { body } = event
  if (!body) {
    return {}
  }
  if (typeof body === 'string') {
    return JSON.parse(body)
  }
  return body
}
