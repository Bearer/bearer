import { BearerState } from '@bearer/core'

class UpdateExistingComponentLifecycleMethods {
  @BearerState() pullRequests: Array<any> = []

  componentWillLoad() {
    console.log('componentWillLoad')
  }

  componentDidUnload() {
    console.log('componentDidUnload')
  }
}
