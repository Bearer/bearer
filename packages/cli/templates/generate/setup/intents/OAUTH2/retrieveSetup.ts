import { TOAUTH2AuthContext, FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/intents'

export default class RetrieveSetupIntent extends FetchData implements FetchData<ReturnedData, any, TOAUTH2AuthContext> {
  async action(event: TFetchActionEvent<Params, TOAUTH2AuthContext>): TFetchPromise<ReturnedData> {
    return { data: event.context.reference }
  }
}

export type Params = {
  setup: any
}

export type ReturnedData = {
  accessToken: string
}
