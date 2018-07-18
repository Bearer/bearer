import axios from 'axios'

import { sendSuccessMessage, sendErrorMessage } from './lambda'

export type TContext = {
  accessToken: string
  [key: string]: any
}

export class Intent {
  static getCollection(
    callback,
    { collection, error }: { collection?: any; error?: any }
  ) {
    if (collection) {
      sendSuccessMessage(callback, collection)
    } else {
      sendErrorMessage(callback, { error: error || 'Unkown error' })
    }
  }

  static getObject(callback, { object, error }: { object?: any; error?: any }) {
    if (object) {
      sendSuccessMessage(callback, object)
    } else {
      sendErrorMessage(callback, { error: error || 'Unkown error' })
    }
  }
}

export const STATE_CLIENT = axios.create({
  timeout: 5000,
  headers: {
    Accept: 'application/json',
    'User-Agent': 'Bearer'
  }
})

class BaseIntent {
  static intent(
    action
  ): (event: any, context: any, callback: (...args: any[]) => any) => any {
    throw new Error(
      'Extending class needs to implement `static intent(action)` method'
    )
  }
}

export class SaveState extends BaseIntent {
  static intent(action) {
    return (event, _context, callback) => {
      const { referenceId } = event.queryStringParameters
      STATE_CLIENT.get(`api/v1/items/${referenceId}`)
        .then(response => {
          console.log('[BEARER]', 'received', response.data)
          const state = response.data.Item
          action(
            event.context,
            event.queryStringParameters,
            event.body,
            state,
            result => {
              STATE_CLIENT.put(`api/v1/items/${referenceId}`, {
                ...result,
                ReadAllowed: true
              })
                .then(data => {
                  console.log('[BEARER]', 'success', data)
                  callback(null, result)
                })
                .catch(e => {
                  console.error('[BEARER]', 'error', e)
                  callback(`Error : ${e}`)
                })
            }
          )
        })
        .catch(response => {
          action(
            event.context,
            event.queryStringParameters,
            event.body,
            {},
            result => {
              STATE_CLIENT.post(`api/v1/items`, {
                ...result,
                ReadAllowed: true
              })
                .then(data => {
                  console.log('[BEARER]', 'success', data)
                  callback(null, result)
                })
                .catch(e => {
                  console.error('[BEARER]', 'error', e)
                  callback(`Error : ${e}`)
                })
            }
          )
        })
    }
  }
}

export class RetrieveState extends BaseIntent {
  static intent(action) {
    return (event, _context, callback) => {
      const { referenceId } = event.queryStringParameters

      STATE_CLIENT.get(`/api/v1/items/${referenceId}`)
        .then(response => {
          if (response.data.error) {
            callback('No data found')
          } else {
            console.log('[BEARER]', 'data', response.data)
            action(
              event.context,
              event.queryStringParameters,
              response.data.Item,
              state => callback(null, state)
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

export class GetCollection extends BaseIntent {
  static intent(action) {
    return (event, _context, callback) =>
      action(event.context, event.queryStringParameters, result => {
        Intent.getCollection(callback, result)
      })
  }
}

export class GetObject extends BaseIntent {
  static intent(action) {
    return (event, _context, callback) =>
      action(event.context, event.queryStringParameters, result => {
        Intent.getObject(callback, result)
      })
  }
}
