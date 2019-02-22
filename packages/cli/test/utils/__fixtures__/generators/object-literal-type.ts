import { FetchData, TFetchPromise, TOAUTH2AuthContext, TFetchActionEvent } from '@bearer/intents'

export default class IntentObjectLiteralType extends FetchData
  implements FetchData<{ expectedData: string[] }, any, TOAUTH2AuthContext> {
  async action(
    event: TFetchActionEvent<
      {
        params: {
          inlineParam: string
          stringEnum: 'none' | 'all' | 'every'
          inlineNumber: number
          nestedObject: { name: string }
        }
      },
      TOAUTH2AuthContext
    >
  ): TFetchPromise<{ expectedData: string[] }> {
    // const token = event.context.authAccess.accessToken
    // Put your logic here
    return { data: { expectedData: [] } }
  }
}
