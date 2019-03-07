// global fetch
import Bearer from './bearer'
import { cleanQuery } from './utils'
import debug from './logger'

const logger = debug.extend('requests')

const defaultInit = {
  headers: {
    'user-agent': 'Bearer',
    'content-type': 'application/json'
  },
  credentials: 'include'
}

export type TBearerRequest<T> = (params: any, init?: any) => Promise<T>

function getClientId(): string {
  const clientId = Bearer.config.clientId
  if (process.env.NODE_ENV !== 'production') {
    logger('No clientId provided, intent')
  }
  return clientId
}

export function bearerRequest<TPromiseReturn>(uri: string, baseParams = {}): TBearerRequest<TPromiseReturn> {
  const url = `${Bearer.config.integrationHost}${uri}`

  return function(params = {}, init = {}): Promise<TPromiseReturn> {
    return new Promise((resolve, reject) => {
      Bearer.instance.maybeInitialized
        .then(() => {
          const sentParams = cleanQuery({
            ...params,
            ...baseParams,
            clientId: getClientId(),
            secured: Bearer.config.secured
          })

          const query = Object.keys(sentParams).map(key => [key, sentParams[key]].join('='))
          const uri = `${url}?${query.join('&')}`
          fetch
            .apply(null, [uri, { ...defaultInit, ...init }])
            .then(async response => {
              const data = await response.json()
              if (response.status > 399) {
                logger('failing request %j', data)
                reject(data)
              } else {
                logger('successful request %j', data)
                resolve(data)
              }
            })
            .catch(e => console.error('Unexpected error 😞', e))
        })
        .catch(() => console.error('[BEARER', 'Error while waiting for authentication'))
    })
  }
}

export function itemRequest<T = any>(): TBearerRequest<T> {
  return bearerRequest('api/v1/items')
}

type TIntentBaseQuery = {
  intentName: string
  scenarioId: string
  setupId: string
}

export function intentRequest<TReturnFormat>({
  intentName,
  scenarioId,
  setupId
}: TIntentBaseQuery): TBearerRequest<TReturnFormat> {
  return bearerRequest(`api/v3/intents/${scenarioId}/${intentName}`, { setupId })
}

export default {
  intentRequest,
  itemRequest
}
