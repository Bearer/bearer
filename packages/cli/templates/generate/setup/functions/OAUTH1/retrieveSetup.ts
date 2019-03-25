import { TOAUTH1AuthContext, FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/functions'

export default class RetrieveSetupFunction extends FetchData implements FetchData<ReturnedData, any, TOAUTH1AuthContext> {
  async action(event: TFetchActionEvent<Params, TOAUTH1AuthContext>): TFetchPromise<ReturnedData> {
    const { data, referenceId } = await event.store.find<ReturnedData>(event.params.referenceId)
    return { data, referenceId }
  }
}

export type Params = {
  setup: any
  referenceId: string
}

export type ReturnedData = {
  accessToken: string
}
