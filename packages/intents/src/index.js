const axios = require('axios')
const { sendSuccessMessage, sendErrorMessage } = require('./lambda')

type TContext = {
  accessToken: string,
  [key: string]: any
}

const Intent = {
  getCollection: (callback, { collection }) => {
    if (collection) {
      sendSuccessMessage(callback, collection)
    } else {
      sendErrorMessage(callback, { error: 'Error' })
    }
  },

  getObject: (callback, { object }) => {
    if (object) {
      sendSuccessMessage(callback, object)
    } else {
      sendErrorMessage(callback, { error: 'Error' })
    }
  }
}

const STATE_CLIENT = axios.create({
  timeout: 5000,
  headers: {
    Accept: 'application/json',
    'User-Agent': 'Bearer'
  }
})

const SaveState = {
  intent(action) {
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
                  callback(null, response.data.Item)
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
                  callback(null, response.data.Item)
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

const RetrieveState = {
  intent(action) {
    return (event, _context, callback) => {
      const { referenceId } = event.queryStringParameters

      STATE_CLIENT.get(`/api/v1/items/${referenceId}`)
        .then(response => {
          if (response.data.error) {
            callback('No data found')
          } else {
            console.log('[BEARER]', 'data', response.data)
            action(
              event.accessToken,
              event.queryStringParameters,
              response.data.Item,
              prs => callback(null, prs)
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

const GetCollection = {
  intent(action) {
    return (event, _context, callback) =>
      action(event.accessToken, event.queryStringParameters, result => {
        Intent.getCollection(callback, result)
      })
  }
}

const GetObject = {
  intent(action) {
    return (event, _context, callback) =>
      action(event.accessToken, event.queryStringParameters, result => {
        Intent.getObject(callback, result)
      })
  }
}

module.exports = {
  Intent,
  GetCollection,
  GetObject,
  SaveState,
  STATE_CLIENT,
  RetrieveState,
  TContext
}
