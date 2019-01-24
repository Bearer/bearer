import { {{contextAuthTypeImport}}FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/intents'

export default class RetrieveSetupIntent extends FetchData implements FetchData<ReturnedData, any{{contextAuthTypeImplements}}> {
  async action(event: TFetchActionEvent<Params{{contextAuthTypeImplements}}>): TFetchPromise<ReturnedData> {
    return { data: event.context.reference }
  }
}

export type Params = {
  setup: any
}

export type ReturnedData = {{dataType}}
