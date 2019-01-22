import { TAPIKEYAuthContext, FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/intents'
// Uncomment this line if you need to use Client
// import Client from './client'

export default class {{intentClassName}}Intent extends FetchData implements FetchData<ReturnedData, any, TAPIKEYAuthContext> {
  async action(event: TFetchActionEvent<Params, TAPIKEYAuthContext>): TFetchPromise<ReturnedData> {
    // const token = event.context.authAccess.apiKey
    // Put your logic here
    return { data: [] }
  }
}

/**
 * Typing
 */
export type Params = {
  // name: string
}

export type ReturnedData = {
  // foo: string[]
}
