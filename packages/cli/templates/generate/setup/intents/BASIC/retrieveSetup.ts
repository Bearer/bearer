import { TBASICAuthContext, FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/intents'

export default class RetrieveSetupIntent extends FetchData implements FetchData<ReturnedData, any, TBASICAuthContext> {
  async action(event: TFetchActionEvent<Params, TBASICAuthContext>): TFetchPromise<ReturnedData> {
    return { data: event.context.reference }
  }
}

export type Params = {
  setup: any
}

export type ReturnedData = {
  username: string
  password: string
}
