/* global fetch */
import Bearer from './Bearer'

export function bearerRequest(uri: string) {
  const url = `${Bearer.config.integrationHost}api/v1/${uri}`

  return function(params = {}, init = {}) {
    const defaultInit = {
      headers: {
        'user-agent': 'Bearer',
        'content-type': 'application/json'
      },
      credentials: 'include'
    }
    return new Promise((resolve, reject) => {
      Bearer.instance.maybeInitialized.then(() => {
        const sentParams = {
          ...params,
          integrationId: Bearer.config.integrationId
        }

        const query = Object.keys(sentParams).map(key =>
          [key, sentParams[key]].join('=')
        )
        const uri = `${url}?${query.join('&')}`
        console.debug('[BEARER]', 'fetch', url, query)
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
    })
  }
}

export function itemRequest() {
  return bearerRequest('items')
}

export function intentRequest({ intentName, scenarioId }) {
  return bearerRequest(`${scenarioId}/${intentName}`)
}

export default {
  intentRequest,
  itemRequest
}
