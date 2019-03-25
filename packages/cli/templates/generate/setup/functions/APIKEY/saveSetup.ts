import { TAPIKEYAuthContext, FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/functions'

export default class SaveSetupFunction extends FetchData implements FetchData<ReturnedData, any, TAPIKEYAuthContext> {
  async action(event: TFetchActionEvent<Params, TAPIKEYAuthContext>): TFetchPromise< ReturnedData> {
    const { data, referenceId } = await event.store.save<State>(event.params.referenceId, event.params.setup)
    return { data, referenceId }
  }
}

export type Params = {
  setup: any
  referenceId?: string
}

export type State = {
  apiKey: string
}


export type ReturnedData = State
