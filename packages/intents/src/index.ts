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

export class SaveState extends StateIntentBase {
  static get display() {
    return 'SaveState'
  }

  static intent(action) {
    return (event, _context, callback) => {
      const { referenceId } = event.queryStringParameters
      console.log('EVENT CONTEXT', event.context)
      const STATE_CLIENT = UserDataClient(event.context.bearerBaseURL)

      try {
        STATE_CLIENT.get(`api/v1/items/${referenceId}`).then(response => {
          console.log('[BEARER]', 'received', response.data)
          const state = response.data.Item
          console.log('STATE', state)
          action(event.context, event.queryStringParameters, event.body, state, result => {
            STATE_CLIENT.put(`api/v1/items/${referenceId}`, {
              ...result,
              ReadAllowed: true
            })
              .then(data => {
                console.log('[BEARER]', 'success', data)
                callback(null, {
                  meta: {
                    referenceId: referenceId
                  },
                  data: {
                    ...result
                  }
                })
              })
              .catch(e => {
                console.error('[BEARER]', 'error', e)
                callback(`Error : ${e}`)
              })
          })
        })
      } catch (e) {
        console.log(e)
        action(event.context, event.queryStringParameters, event.body, {}, result => {
          STATE_CLIENT.post(`api/v1/items`, {
            ...result,
            ReadAllowed: true
          })
            .then((response: d.TStateData) => {
              console.log('[BEARER]', 'success', response.data)
              callback(null, {
                meta: {
                  referenceId: response.data.Item.referenceId
                },
                data: {
                  ...result
                }
              })
            })
            .catch(e => {
              console.error('[BEARER]', 'error', e)
              callback(`Error : ${e}`)
            })
        })
      }
    }
  }
}

export class RetrieveState extends StateIntentBase {
  static get display() {
    return 'RetrieveState'
  }

  static intent(action) {
    return (event, context, callback) => {
      const { referenceId } = event.queryStringParameters
      const STATE_CLIENT = UserDataClient(event.context.bearerBaseURL)

      STATE_CLIENT.get(`/api/v1/items/${referenceId}`)
        .then(response => {
          if (response.data.error) {
            callback('No data found')
          } else {
            console.log('[BEARER]', 'data', response.data)
            action(event.context, event.queryStringParameters, response.data.Item, state =>
              callback(null, {
                meta: {
                  referenceId: response.data.Item.referenceId
                },
                data: state
              })
            )
          }
        })
        .catch(e => {
          callback('No data found')
          console.log('[BEARER]', 'error', e)
        })
    }
  }
}

export class GetCollection extends GenericIntentBase {
  static get display() {
    return 'GetCollection'
  }

  static intent(action) {
    return (event, _context, callback) =>
      action(event.context, event.queryStringParameters, result => {
        Intent.getCollection(callback, result)
      })
  }
}

export class GetObject extends GenericIntentBase {
  static get display() {
    return 'GetObject'
  }

  static intent(action) {
    return (event, _context, callback) =>
      action(event.context, event.queryStringParameters, result => {
        Intent.getObject(callback, result)
      })
  }
}
