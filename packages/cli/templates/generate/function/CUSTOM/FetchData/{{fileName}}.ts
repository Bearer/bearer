import { TCUSTOMAuthContext, FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/functions'
// Uncomment this line if you need to use Client
// import Client from './client'

export default class {{functionClassName}}Function extends FetchData implements FetchData<ReturnedData, any, TCUSTOMAuthContext> {
  async action(event: TFetchActionEvent<Params, TCUSTOMAuthContext>): TFetchPromise<ReturnedData> {
    // Put your logic here
    return { data: [] }
  }

  // Uncomment the line below if you don't want your function to be called from the frontend
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
