import axios, { AxiosInstance } from 'axios'
import * as d from './declaration'
export * from './declaration'

import { sendSuccessMessage, sendErrorMessage } from './lambda'
import UserDataClient from './UserDataClient'

export class Intent {
  static fetchData(callback, { data, error }: { data?: any; error?: any }) {
    if (data) {
      sendSuccessMessage(callback, { data })
    } else {
      sendErrorMessage(callback, { error: error || 'Unkown error' })
    }
  }
}

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

type ISaveStateIntentAction = {
  (context: any, _params: any, body: any, state: any, callback: (state: any) => void): void
}

export class SaveState extends StateIntentBase {
  static get display() {
    return 'SaveState'
  }

  static intent(action: ISaveStateIntentAction) {
    return (event, _context, lambdaCallback) => {
      const { referenceId } = event.queryStringParameters
      const STATE_CLIENT = UserDataClient(event.context.bearerBaseURL)
      try {
        STATE_CLIENT.retrieveState(referenceId)
          .then(savedState => {
            const state = savedState ? savedState.Item : {}
            try {
              action(event.context, event.queryStringParameters, event.body, state, result => {
                if (savedState) {
                  STATE_CLIENT.updateState(referenceId, result)
                    .then(() => {
                      lambdaCallback(null, { meta: { referenceId }, data: { ...result } })
                    })
                    .catch(error => lambdaCallback(error.toString(), { error: error.toString() }))
                } else {
                  STATE_CLIENT.saveState(result)
                    .then(data => {
                      lambdaCallback(null, { meta: { referenceId: data.Item.referenceId }, data: { ...result } })
                    })
                    .catch(error => lambdaCallback(error.toString(), { error: error.toString() }))
                }
              })
            } catch (error) {
              return lambdaCallback(error.toString(), { error: error.toString() })
            }
          })
          .catch(error => lambdaCallback(error.toString(), { error: error.toString() }))
      } catch (error) {
        return lambdaCallback(error.toString(), { error: error.toString() })
      }
    }
  }
}

export class RetrieveState extends StateIntentBase {
  static get display() {
    return 'RetrieveState'
  }

  static intent(action) {
    return (event, _context, lambdaCallback) => {
      const { referenceId } = event.queryStringParameters
      const STATE_CLIENT = UserDataClient(event.context.bearerBaseURL)
      try {
        STATE_CLIENT.retrieveState(referenceId).then(state => {
          if (state) {
            action(event.context, event.queryStringParameters, state.Item, preparedState => {
              lambdaCallback(null, {
                meta: {
                  referenceId: state.Item.referenceId
                },
                data: preparedState
              })
            })
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
  static get display() {
    return 'FetchData'
  }

  static intent(action) {
    return (event, _context, lambdaCallback) =>
      action(event.context, event.queryStringParameters, result => {
        Intent.fetchData(lambdaCallback, result)
      })
  }
}
