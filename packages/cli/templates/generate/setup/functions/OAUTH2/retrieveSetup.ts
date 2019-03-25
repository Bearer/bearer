import { TOAUTH2AuthContext, FetchData, TFetchActionEvent, TFetchPromise } from '@bearer/functions'

export default class RetrieveSetupFunction extends FetchData implements FetchData<ReturnedData, any, TOAUTH2AuthContext> {
  async action(event: TFetchActionEvent<Params, TOAUTH2AuthContext>): TFetchPromise<ReturnedData> {
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
  referenceId: string
}
