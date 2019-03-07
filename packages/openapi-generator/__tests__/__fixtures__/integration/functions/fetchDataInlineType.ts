import { FetchData, TFetchActionEvent, TFetchPromise, TOAUTH2AuthContext } from '@bearer/functions'

import { PullRequest } from '../views/types'

export default class FetchDataInlineType extends FetchData
  implements FetchData<PullRequest[], any, TOAUTH2AuthContext> {
  async action(event: TFetchActionEvent<{ pullRequests: string[] }, TOAUTH2AuthContext>): TFetchPromise<PullRequest[]> {
    return { data: [] }
  }
}
