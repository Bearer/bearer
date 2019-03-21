import { TOAUTH1AuthContext, FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/functions'

export default class SaveSetupFunction extends FetchData implements FetchData<ReturnedData, any, TOAUTH1AuthContext> {
  async action(event: TFetchActionEvent<Params, TOAUTH1AuthContext>): TFetchPromise< ReturnedData> {
    const { data, referenceId } = await event.store.save<State>(event.params.referenceId, event.params.setup)
    return { data,  referenceId }
  }
}

export type Params = {
  setup: any
  referenceId?: string
}

export type State = {
  consumerKey: string
  consumerSecret: string
}

export type ReturnedData = State
