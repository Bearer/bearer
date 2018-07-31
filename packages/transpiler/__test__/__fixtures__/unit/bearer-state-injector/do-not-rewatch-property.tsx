import { BearerState, Watch } from '@bearer/core'

class UpdateExistingPropertyWatcher {
  @BearerState() pullRequests: Array<any> = []

  @Watch('pullRequests')
  pullRequestsChangeHandler(newValue, oldValue) {
    console.log('Prepend stuff')
  }
}
