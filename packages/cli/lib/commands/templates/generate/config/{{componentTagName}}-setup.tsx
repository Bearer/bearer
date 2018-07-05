/*
  The purpose of this component is to save scenario credentials.
  This file has been generated automatically and should not be edited.
*/

import { Component, State } from '@bearer/core'
import '@bearer/ui'

@Component({
  tag: '{{componentTagName}}-config',
  shadow: true
})
export class {{scenarioTitle}}Config {
  @State() fields = {{fields}}
  @State() innerListener = `config_success:BEARER_SCENARIO_ID`
  render() {
    return (
      <div>
        <bearer-dropdown-button innerListener={this.innerListener}>
          <span slot="buttonText">Config component</span>
          <bearer-config scenarioId="BEARER_SCENARIO_ID" fields={this.fields} />
        </bearer-dropdown-button>
      </div>
    )
  }
}

