import { TBASICAuthContext, FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/intents'
// Uncomment this line if you need to use Client
// import Client from './client'

export default class {{intentClassName}}Intent extends FetchData implements FetchData<ReturnedData, any, TBASICAuthContext>{
  async action(event: TFetchActionEvent<Params, TBASICAuthContext>): TFetchPromise<ReturnedData> {
    // const { username, password } = event.context.authAccess
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
