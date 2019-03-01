import { SaveState, TOAUTH1AuthContext, TSaveActionEvent, TSavePromise } from '@bearer/intents'

import { PullRequest } from '../views/types'

export default class SavePullRequestsIntent extends SaveState
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
