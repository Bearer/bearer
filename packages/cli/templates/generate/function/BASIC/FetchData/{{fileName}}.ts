import { TBASICAuthContext, FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/functions'
// Uncomment this line if you need to use Client
// import Client from './client'

export default class {{functionClassName}}Function extends FetchData implements FetchData<ReturnedData, any, TBASICAuthContext>{
  async action(event: TFetchActionEvent<Params, TBASICAuthContext>): TFetchPromise<ReturnedData> {
    // const { username, password } = event.context.authAccess
    // Put your logic here
    return { data: [] }
  }

  // Uncomment the line above if you don't want your function to be called from the frontend
  // static backendOnly = true
  
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
