import { TAPIKEYAuthContext, FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/intents'

export default class RetrieveSetupIntent extends FetchData implements FetchData<ReturnedData, any, TAPIKEYAuthContext> {
  async action(event: TFetchActionEvent<Params, TAPIKEYAuthContext>): TFetchPromise<ReturnedData> {
    return { data: event.context.reference }
  }
}

export type Params = {
  setup: any
}

export type ReturnedData = {
  apiKey: string
}
