import { Component, Prop, State } from '@bearer/core'

import Bearer from '@bearer/core'

@Component({
  tag: 'bearer-setup-display',
  shadow: true
})
export class BearerSetupDisplay {
  @Prop() setupId = ''
  @State() isSetup: boolean = false

  componentDidLoad() {
    Bearer.emitter.addListener(`setup_success:${this.setupId}`, () => {
      this.isSetup = true
    })
  }

  render() {
    if (this.isSetup) {
      return (
        <p>
          Scenario is currently setup with Setup ID:&nbsp;
          <bearer-badge kind="info">{this.setupId}</bearer-badge>
        </p>
      )
    } else {
      return (
        <p>
          <p>Scenario hasn't been setup yet</p>
        </p>
      )
    }
  }
}
