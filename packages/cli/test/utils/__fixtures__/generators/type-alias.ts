import { TOAUTH2AuthContext, TFetchPayload, FetchData, TFetchActionEvent } from '@bearer/intents'
// Uncomment this line if you need to use Client
// import Client from './client'

export default class Intent {
  static intentType = FetchData

  static async action(event: TFetchActionEvent<TOAUTH2AuthContext, Params>): Promise<ReturnedData> {
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

export type ReturnedData = TFetchPayload<{
  foo: string[]
  anObject: {
    values?: number[]
  }
}>
