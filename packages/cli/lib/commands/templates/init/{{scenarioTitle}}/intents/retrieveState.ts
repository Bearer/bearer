import { RetrieveState } from './intents'
import { ScenarioState, SavedData, PullRequest } from './saveState'

interface RetrieveStateCallback {
  ({ pullRequests }: { pullRequests: Array<PullRequest> }): void
}

interface IAction {
  (
    token: string,
    _params,
    state: ScenarioState,
    callback: RetrieveStateCallback
  ): void
}

// type IPromises = Array<Promise<PullRequest>>

export default class RetrieveStateIntent {
  static get intentName(): string {
    return 'retrieveState'
  }
  static get intentType() {
    return RetrieveState
  }

  static action: IAction = (token, _params, state, callback) => {
    // const pullRequestsPromises: IPromises = state.pullRequests.map(
    //   (pr: SavedData) =>
    //     new Promise((resolve, _reject) => {
    //       GetPullHelloWorldsIntent.action(
    //         token,
    //         { id: pr.number, fullName: pr.full_name },
    //         ({ object }: { object: PullRequest }) => resolve(object)
    //       )
    //     })
    // )
    // Promise.all(pullRequestsPromises).then(prs => {
    //   console.log('[BEARER]', 'prs', prs)
    //   callback({ pullRequests: prs.map(pr => pr) })
    // })
    callback(state)
  }
}
