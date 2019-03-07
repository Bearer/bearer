import { FetchData, TFetchActionEvent, TFetchPromise, TOAUTH2AuthContext } from '@bearer/intents'

import { Repo } from '../views/types'

export default class FetchDataRefType extends FetchData implements FetchData<ReturnedData, any, TOAUTH2AuthContext> {
  async action(event: TFetchActionEvent<Params, TOAUTH2AuthContext>): TFetchPromise<ReturnedData> {
    return {
      data: []
    }
  }
}
export type Params = {
  reference: string
}

export type ReturnedData = Repo[]
