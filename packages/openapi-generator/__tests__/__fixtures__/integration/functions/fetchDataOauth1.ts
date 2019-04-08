import { FetchData, TFetchActionEvent, TFetchPromise, TOAUTH1AuthContext } from '@bearer/functions'

import { PullRequest } from '../views/types'

export default class SearchPullRequestsFunction extends FetchData
  implements FetchData<ReturnedData, any, TOAUTH1AuthContext> {
  async action(event: TFetchActionEvent<Params, TOAUTH1AuthContext>): TFetchPromise<ReturnedData> {
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
