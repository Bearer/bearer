import { FetchData, TFetchActionEvent, TFetchPromise, TAPIKEYAuthContext } from '@bearer/functions'

import { PullRequest } from '../views/types'

export default class SearchPullRequestsFunction extends FetchData
  implements FetchData<ReturnedData, any, TAPIKEYAuthContext> {
  async action(event: TFetchActionEvent<Params, TAPIKEYAuthContext>): TFetchPromise<ReturnedData> {
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
