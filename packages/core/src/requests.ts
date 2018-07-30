/* global fetch */
import Bearer from './Bearer'

const defaultInit = {
  headers: {
    'user-agent': 'Bearer',
    'content-type': 'application/json'
  },
  credentials: 'include'
}

export function bearerRequest(uri: string, baseParams = {}): (params: any, init?: any) => Promise<any> {
  const url = `${Bearer.config.integrationHost}api/v1/${uri}`

  return function(params = {}, init = {}): Promise<any> {
    return new Promise((resolve, reject) => {
      Bearer.instance.maybeInitialized
        .then(() => {
          const sentParams = {
            ...params,
            ...baseParams,
            integrationId: Bearer.config.integrationId
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

export function itemRequest() {
  return bearerRequest('items')
}

export function intentRequest({
  intentName,
  scenarioId,
  setupId
}: {
  intentName: string
  scenarioId: string
  setupId: string
}) {
  return bearerRequest(`${scenarioId}/${intentName}`, { setupId })
}

export default {
  intentRequest,
  itemRequest
}
