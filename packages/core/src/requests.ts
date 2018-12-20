// global fetch
import Bearer from './bearer'

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
    console.debug('[BEARER]', 'No clientId provided, intent')
  }
  return clientId
}

export function bearerRequest<TPromiseReturn>(uri: string, baseParams = {}): TBearerRequest<TPromiseReturn> {
  const url = `${Bearer.config.integrationHost}api/v2/intents/${uri}`

  return function(params = {}, init = {}): Promise<TPromiseReturn> {
    return new Promise((resolve, reject) => {
      Bearer.instance.maybeInitialized
        .then(() => {
          const sentParams = {
            ...params,
            ...baseParams,
            clientId: getClientId(),
            secured: Bearer.config.secured
          }

          const query = Object.keys(sentParams).map(key => [key, sentParams[key]].join('='))
          const uri = `${url}?${query.join('&')}`
          fetch
            .apply(null, [uri, { ...defaultInit, ...init }])
            .then(response => {
              const data = response.json()
              if (response.status > 399) {
                console.debug('[BEARER]', 'failing request', response)
                reject(data)
              } else {
                console.debug('[BEARER]', 'successful request', response)
                resolve(data)
              }
            })
            .catch(e => console.error('Unexpected error 😞', e))
        })
        .catch(() => console.error('[BEARER', 'Error while waiting for authentication'))
    })
  }
}

export function itemRequest(): TBearerRequest<any> {
  return bearerRequest('items')
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
  return bearerRequest(`${scenarioId}/${intentName}`, { setupId })
}

export default {
  intentRequest,
  itemRequest
}
