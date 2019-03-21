import { TBASICAuthContext, FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/functions'

export default class SaveSetupFunction extends FetchData implements FetchData<ReturnedData, any, TBASICAuthContext> {
  async action(event: TFetchActionEvent<Params, TBASICAuthContext>): TFetchPromise< ReturnedData> {
    const { data, referenceId } = await event.store.save<State>(event.params.referenceId, event.params.setup)
    return { data,  referenceId }
  }
}

export type Params = {
  setup: any
  referenceId?: string
}

export type State = {
  username: string
  password: string
}

export type ReturnedData = State
