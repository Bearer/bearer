import { TOAUTH2AuthContext, TFetchPayload, FetchData, TFetchActionEvent } from '@bearer/intents'
// Uncomment this line if you need to use Client
// import Client from './client'

export default class Intent {
  static intentType = FetchData

  static async action(event: TFetchActionEvent<TOAUTH2AuthContext, Params>): Promise<ReturnedData> {
    // const token = event.context.authAccess.accessToken
    // Put your logic here
    return { data: [] }
  }
}

/**
 * Typing
 */
export type Params = {
  name: string
  somethingElse: number
}

export type ReturnedData = TFetchPayload<{
  // foo: string[]
}>
