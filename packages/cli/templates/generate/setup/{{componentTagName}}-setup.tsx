/*
  The purpose of this component is to save scenario credentials.
  This file has been generated automatically and should not be edited.
*/

import { Component, State, Prop } from '@bearer/core'
import '@bearer/ui'

@Component({
  tag: '{{componentTagName}}-setup',
  shadow: true
})
export class {{scenarioTitle}}Setup {
  @Prop() onSetupSuccess: (detail: any) => void = (_any: any) => {}
  @State() fields = {{fields}}
  @State() innerListener = `setup_success:BEARER_SCENARIO_ID`
  render() {
    return (
      <bearer-dropdown-button innerListener={this.innerListener} btnProps={ { content: "Setup component" } }>
        <bearer-setup onSetupSuccess={this.onSetupSuccess} scenarioId="BEARER_SCENARIO_ID" fields={this.fields} />
      </bearer-dropdown-button>
    )
  }
}

