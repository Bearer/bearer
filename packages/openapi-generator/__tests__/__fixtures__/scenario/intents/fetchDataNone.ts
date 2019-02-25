import { FetchData, TFetchActionEvent, TFetchPromise, TNONEAuthContext } from '@bearer/intents'

import { PullRequest } from '../views/types'

export default class SearchPullRequestsIntent extends FetchData
  implements FetchData<ReturnedData, any, TNONEAuthContext> {
  async action(event: TFetchActionEvent<Params, TNONEAuthContext>): TFetchPromise<ReturnedData> {
    return { data: [] }
  }
}
type ParamsId = {
  id: string
}

type ParamsSearch = {
  name: string
  query: string
}
export type ReturnedData = PullRequest[]
export type Params = ParamsId | ParamsSearch
