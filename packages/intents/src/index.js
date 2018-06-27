const { sendSuccessMessage, sendErrorMessage } = require('./lambda')

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

const StoreObject = {
  intent(action) {
    return (event, context, callback) =>
      action(
        event.accessToken,
        event.queryStringParameters,
        event.body,
        result => {
          Intent.getObject(callback, result)
        }
      )
  }
}

const GetCollection = {
  intent(action) {
    return (event, context, callback) =>
      action(event.accessToken, event.queryStringParameters, result => {
        Intent.getCollection(callback, result)
      })
  }
}

const GetObject = {
  intent(action) {
    return (event, context, callback) =>
      action(event.accessToken, event.queryStringParameters, result => {
        Intent.getObject(callback, result)
      })
  }
}

const SaveState = {
  intent(action) {
    return (event, _context, callback) => {
      const { referenceId } = event.queryStringParameters
      client
        .get(`api/v1/items/${referenceId}`)
        .then(response => {
          console.log('[BEARER]', 'received', response.data)
          const state = response.data.Item
          action(
            event.accessToken,
            event.queryStringParameters,
            event.body,
            state,
            result => {
              client
                .put(`api/v1/items/${referenceId}`, {
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
            event.accessToken,
            event.queryStringParameters,
            event.body,
            {},
            result => {
              client
                .post(`api/v1/items`, {
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

module.exports = { Intent, GetCollection, GetObject, StoreObject, SaveState }
