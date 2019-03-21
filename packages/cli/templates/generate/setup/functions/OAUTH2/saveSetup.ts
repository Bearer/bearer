import { TOAUTH2AuthContext, FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/functions'

export default class SaveSetupFunction extends FetchData implements FetchData<ReturnedData, any, TOAUTH2AuthContext> {
  async action(event: TFetchActionEvent<Params, TOAUTH2AuthContext>): TFetchPromise< ReturnedData> {
    const { data, referenceId } = await event.store.save<State>(event.params.referenceId, event.params.setup)
    return { data,  referenceId }
  }
}

export type Params = {
  setup: any
  referenceId?: string
}

export type State = {
  accessToken: string
}

export type ReturnedData = State
