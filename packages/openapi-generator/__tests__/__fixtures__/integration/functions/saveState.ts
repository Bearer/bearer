import { SaveState, TOAUTH1AuthContext, TSaveActionEvent, TSavePromise } from '@bearer/functions'

import { PullRequest } from '../views/types'

export default class SavePullRequestsFunction extends SaveState
  implements SaveState<State, ReturnedData, any, TOAUTH1AuthContext> {
  async action(event: TSaveActionEvent<State, Params, TOAUTH1AuthContext>): TSavePromise<State, ReturnedData> {
    return {
      state: {
        pullRequests: []
      },
      data: []
    }
  }
}

/**
 * Typing
 */
export type Params = {
  pullRequests: PullRequest[]
}

export type State = {
  pullRequests: string[]
}

export type ReturnedData = PullRequest[]
