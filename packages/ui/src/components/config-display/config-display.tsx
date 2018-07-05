import { Component, Prop, State } from '@bearer/core'

import Bearer from '@bearer/core'

@Component({
  tag: 'bearer-config-display',
  shadow: true
})
export class BearerConfigDisplay {
  @Prop() scenarioId = ''
  @State() isConfig: boolean = false
  @State() configId = ''

  componentDidLoad() {
    Bearer.emitter.addListener(`config_success:${this.scenarioId}`, data => {
      this.configId = data.referenceID
      this.isConfig = true
    })
  }

  render() {
    if (this.isConfig) {
      return (
        <p>
          Scenario is currently configure with Config ID:&nbsp;
          <bearer-badge kind="info">{this.configId}</bearer-badge>
        </p>
      )
    } else {
      return (
        <p>
          <p>Scenario hasn't been configured yet</p>
        </p>
      )
    }
  }
}
