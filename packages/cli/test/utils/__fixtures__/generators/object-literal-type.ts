import { FetchData, TFetchPromise } from '@bearer/intents'
// Uncomment this line if you need to use Client
// import Client from './client'

export default class Intent extends FetchData implements FetchData {
  async action(event: {
    params: {
      inlineParam: string
      stringEnum: 'none' | 'all' | 'every'
      inlineNumber: number
      nestedObject: { name: string }
    }
  }): TFetchPromise<{ expectedData: string[] }> {
    // const token = event.context.authAccess.accessToken
    // Put your logic here
    return { data: { expectedData: [] } }
  }
}
