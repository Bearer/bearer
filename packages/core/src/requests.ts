// global fetch
import Bearer from './Bearer'

const defaultInit = {
  headers: {
    'user-agent': 'Bearer',
    'content-type': 'application/json'
  },
  credentials: 'include'
}

export type TBearerRequest<T> = (params: any, init?: any) => Promise<T>

export function bearerRequest<TPromiseReturn>(uri: string, baseParams = {}): TBearerRequest<TPromiseReturn> {
  const url = `${Bearer.config.integrationHost}api/v1/${uri}`

  return function(params = {}, init = {}): Promise<TPromiseReturn> {
    return new Promise((resolve, reject) => {
      Bearer.instance.maybeInitialized
        .then(() => {
          const sentParams = {
            ...params,
            ...baseParams,
            clientId: Bearer.config.clientId
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
            .catch(e => console.log('Unexpected error 😞', e))
        })
        .catch(() => console.error('[BEARER', 'Erro while waiting for authentication'))
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
