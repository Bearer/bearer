import { BearerState, Watch } from '@bearer/core'

class UpdateExistingPropertyWatcher {
  @BearerState()
  pullRequests: Array<any> = []

  @BearerState({ statePropName: 'repo' })
  repository: {} = {}

  @Watch('pullRequests')
  pullRequestsChangeHandler(newPullRequest: Array<any>, oldValue: any[]) {
    console.log('Prepend stuff')
  }
}
