import axios, { AxiosInstance } from 'axios'
import * as d from './declaration'
export * from './declaration'

import { sendSuccessMessage, sendErrorMessage } from './lambda'
import UserDataClient from './UserDataClient'

export class Intent {
  static getCollection(callback, { collection, error }: { collection?: any; error?: any }) {
    if (collection) {
      sendSuccessMessage(callback, { data: collection })
    } else {
      sendErrorMessage(callback, { error: error || 'Unkown error' })
    }
  }

  static getObject(callback, { object, error }: { object?: any; error?: any }) {
    if (object) {
      sendSuccessMessage(callback, { data: object })
    } else {
      sendErrorMessage(callback, { error: error || 'Unkown error' })
    }
  }
}

class BaseIntent {
  static get display(): string {
    throw new Error('Extending class needs to implement `static intent(action)` method')
  }

  static intent(action): (event: any, context: any, callback: (...args: any[]) => any) => any {
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

export class GetCollection extends GenericIntentBase {
  static get display() {
    return 'GetCollection'
  }

  static intent(action) {
    return (event, _context, lambdaCallback) =>
      action(event.context, event.queryStringParameters, result => {
        Intent.getCollection(lambdaCallback, result)
      })
  }
}

export class GetObject extends GenericIntentBase {
  static get display() {
    return 'GetObject'
  }

  static intent(action) {
    return (event, _context, lambdaCallback) =>
      action(event.context, event.queryStringParameters, result => {
        Intent.getObject(lambdaCallback, result)
      })
  }
}
