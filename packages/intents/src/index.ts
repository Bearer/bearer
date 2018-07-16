import axios from 'axios'

import { sendSuccessMessage, sendErrorMessage } from './lambda'

export type TContext = {
  accessToken: string
  [key: string]: any
}

export class Intent {
  static getCollection(callback, { collection }) {
    if (collection) {
      sendSuccessMessage(callback, collection)
    } else {
      sendErrorMessage(callback, { error: 'Error' })
    }
  }

  static getObject(callback, { object }) {
    if (object) {
      sendSuccessMessage(callback, object)
    } else {
      sendErrorMessage(callback, { error: 'Error' })
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
  static get template(): string {
    throw new Error(
      'Extending class needs to implement `static get template()` method'
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

export class SaveState extends BaseIntent {
  static get template() {
    return `
  STATE_CLIENT.defaults.baseURL = 'https://int.bearer.sh/'

  static action(
    _context,
    _params,
    body: any,
    state: any,
    callback: (any) => void
  ): void {
    const { item: { name } } = body
    const { items = [] }: any = state
    const newItem: any = { name }

    callback({
      ...state,
      items: [...items, newItem]
    })
  }
`
  }
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

export class RetrieveState extends BaseIntent {
  static get template(): string {
    return `
  STATE_CLIENT.defaults.baseURL = 'https://int.bearer.sh/'

  static action(_context: TContext, _params: any, state, callback) => {
    callback({ items: state.items.map(({name}) => name) })
  }
`
  }
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

export class GetCollection extends BaseIntent {
  static get template(): string {
    return `
  static action(context: TContext, params: any, callback: (params: any) => void) {
    //... your code goes here
    // use the client defined in client.ts to fetch real object like that:
    // CLIENT.get('/people', { headers: headersFor(context.accessToken) }).then(({ data }) => {
    //   callback({ collection: data.results });
    // });
    callback({ collection: []})
  }
`
  }
  static intent(action) {
    return (event, _context, callback) =>
      action(event.context, event.queryStringParameters, result => {
        Intent.getCollection(callback, result)
      })
  }
}

export class GetObject extends BaseIntent {
  static get template(): string {
    return `
  static action(context: TContext, params: any, callback: (params: any) => void) {
    //... your code goes here
    // use the client defined in client.ts to fetch real object like that:
    // CLIENT.get(\`/people/\${params.id}\`, { headers: headersFor(context.accessToken) })
    //   .then(({ data }) => {
    //     callback({ object: data });
    //   });
    callback({ object: {}})
  }
`
  }
  static intent(action) {
    return (event, _context, callback) =>
      action(event.context, event.queryStringParameters, result => {
        Intent.getObject(callback, result)
      })
  }
}
