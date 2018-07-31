import { BearerState } from '@bearer/core'

class UpdateExistingComponentLifecycleMethods {
  @BearerState() pullRequests: Array<any> = []

  // TODO: handle these use cases

  componentWillLoad() {
    console.log('componentWillLoad')
  }

  componentDidUnload() {
    console.log('componentDidUnload')
  }
}
