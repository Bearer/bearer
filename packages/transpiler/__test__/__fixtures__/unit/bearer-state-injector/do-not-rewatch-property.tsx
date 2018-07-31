import { BearerState, Watch } from '@bearer/core'

class UpdateExistingPropertyWatcher {
  @BearerState() pullRequests: Array<any> = []

  // TODO: handle this use case
  @Watch('pullRequests')
  pullRequestsChangeHandler(newValue, oldValue) {
    console.log('Prepend stuff')
  }
}
