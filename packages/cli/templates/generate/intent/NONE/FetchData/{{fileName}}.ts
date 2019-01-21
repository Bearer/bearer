import { TNONEAuthContext, TFetchPayload, FetchData, TFetchActionEvent } from '@bearer/intents'
// Uncomment this line if you need to use Client
// import Client from './client'

export default class {{intentClassName}}Intent {
  intentType = FetchData

  async action(event: TFetchActionEvent<TNONEAuthContext, Params>): Promise<ReturnedData> {
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

export type ReturnedData = TFetchPayload<{
  // foo: string[]
}>
