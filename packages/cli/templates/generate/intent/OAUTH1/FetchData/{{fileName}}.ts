import { TOAUTH1AuthContext, FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/intents'
// Uncomment this line if you need to use Client
// import Client from './client'

export default class {{intentClassName}}Intent extends FetchData implements FetchData<ReturnedData, any, TOAUTH1AuthContext> {
  async action(event: TFetchActionEvent<Params, TOAUTH1AuthContext>): TFetchPromise<ReturnedData> {
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
