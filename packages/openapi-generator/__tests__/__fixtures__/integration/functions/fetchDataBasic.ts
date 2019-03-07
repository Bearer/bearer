import { FetchData, TFetchActionEvent, TFetchPromise, TBASICAuthContext } from '@bearer/functions'

import { PullRequest } from '../views/types'

export default class SearchPullRequestsFunction extends FetchData
  implements FetchData<ReturnedData, any, TBASICAuthContext> {
  async action(event: TFetchActionEvent<Params, TBASICAuthContext>): TFetchPromise<ReturnedData> {
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
