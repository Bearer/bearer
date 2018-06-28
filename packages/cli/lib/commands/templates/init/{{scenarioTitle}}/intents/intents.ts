import axios from 'axios'
const client = axios.create({
  baseURL: 'https://int.dev.bearer.sh/',
  timeout: 5000,
  headers: {
    Accept: 'application/json',
    'User-Agent': 'Bearer'
  }
})

interface ICallback {
  (params: null | string, payload?: Object | null): void
}

interface IDynamoData {
  data: {
    Item: {
      referenceId: string
      [name: string]: any
    }
  }
}

interface ILambdaEvent {
  accessToken: string
  queryStringParameters: {
    [key: string]: any
  }
  body: any
}

interface ILambda {
  (event: ILambdaEvent, _context: any, callback: ICallback): void
}

interface ISaveAction {
  (
    token: string,
    params: any,
    body: any,
    state: any,
    callback: (state: any) => void
  ): void
}

export const SaveState = {
  intent(action: ISaveAction): ILambda {
    return (event, _context, callback: ICallback) => {
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
                .then((data: IDynamoData) => {
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
                .then((data: IDynamoData) => {
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

interface IRetrieveAction {
  (token: string, params: any, state: any, callback: (state: any) => void): void
}

export const RetrieveState = {
  intent(action: IRetrieveAction): ILambda {
    return (event: ILambdaEvent, _context, callback: ICallback) => {
      const { referenceId } = event.queryStringParameters

      client
        .get(`/api/v1/items/${referenceId}`)
        .then((response: any) => {
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
