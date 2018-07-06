import { Component, Prop, State } from '@bearer/core'

import Bearer from '@bearer/core'

@Component({
  tag: 'bearer-setup-display',
  shadow: true
})
export class BearerSetupDisplay {
  @Prop() scenarioId = ''
  @State() isSetup: boolean = false
  @State() setupId = ''

  componentDidLoad() {
    Bearer.emitter.addListener(`setup_success:${this.scenarioId}`, data => {
      this.setupId = data.referenceId
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
