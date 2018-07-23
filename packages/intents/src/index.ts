import axios, { AxiosInstance, AxiosResponse } from 'axios'

import { sendSuccessMessage, sendErrorMessage } from './lambda'

export type TContext = {
  accessToken: string
  [key: string]: any
}

export type TStateData = AxiosResponse<{
  Item: any
}>

export class Intent {
  static getCollection(
    callback,
    { collection, error }: { collection?: any; error?: any }
  ) {
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

export const STATE_CLIENT: AxiosInstance = axios.create({
  baseURL: 'https://int.staging.bearer.sh',
  timeout: 5000,
  headers: {
    Accept: 'application/json',
    'User-Agent': 'Bearer'
  }
})

class BaseIntent {
  static get display(): string {
    throw new Error(
      'Extending class needs to implement `static intent(action)` method'
    )
  }

  static intent(
    action
  ): (event: any, context: any, callback: (...args: any[]) => any) => any {
    throw new Error(
      'Extending class needs to implement `static intent(action)` method'
    )
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
    return async (event, context, callback) => {
      const { referenceId } = event.queryStringParameters
      const baseURL =
        event.context.bearerBaseURL || STATE_CLIENT.defaults.baseURL
      try {
        const response = await STATE_CLIENT.request({
          method: 'get',
          url: `api/v1/items/${referenceId}`,
          baseURL
        })
        console.log('[BEARER]', 'received', response.data)
        const state = response.data.Item
        action(
          event.context,
          event.queryStringParameters,
          event.body,
          state,
          result => {
            STATE_CLIENT.request({
              method: 'put',
              url: `api/v1/items/${referenceId}`,
              baseURL,
              data: {
                ...result,
                ReadAllowed: true
              }
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
          }
        )
      } catch (e) {
        console.log(e)
        action(
          event.context,
          event.queryStringParameters,
          event.body,
          {},
          result => {
            STATE_CLIENT.request({
              method: 'post',
              url: `api/v1/items`,
              baseURL,
              data: {
                ...result,
                ReadAllowed: true
              }
            })
              .then((response: TStateData) => {
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
          }
        )
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
      const baseURL =
        event.context.bearerBaseURL || STATE_CLIENT.defaults.baseURL

      STATE_CLIENT.request({
        method: 'get',
        url: `/api/v1/items/${referenceId}`,
        baseURL
      })
        .then(response => {
          if (response.data.error) {
            callback('No data found')
          } else {
            console.log('[BEARER]', 'data', response.data)
            action(
              event.context,
              event.queryStringParameters,
              response.data.Item,
              state =>
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
