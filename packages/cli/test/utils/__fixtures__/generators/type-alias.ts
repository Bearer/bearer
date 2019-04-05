import { TOAUTH2AuthContext, FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/functions'
// Uncomment the line below to use the API Client
// import Client from './client'

export default class FunctionTypeAlias extends FetchData implements FetchData<ReturnedData, any, TOAUTH2AuthContext> {
  action = async (event: TFetchActionEvent<Params, TOAUTH2AuthContext>): TFetchPromise<ReturnedData> => {
    // const token = event.context.authAccess.accessToken
    // Put your logic here
    return {
      data: {
        foo: ['all', 'none'],
        anObject: {
          values: [1, 2]
        }
      }
    }
  }
}

/**
 * Typing
 */
export type Params = {
  aliasParam: string
  stringEnum: 'none' | 'all' | 'every'
  inlineNumber: number
  nestedObject: { name: string }
}

export type ReturnedData = {
  foo: string[]
  anObject: {
    values?: number[]
  }
}
