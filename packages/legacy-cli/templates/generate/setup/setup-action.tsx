/*
  The purpose of this component is to save scenario credentials.
  This file has been generated automatically and should not be edited.
*/

import { RootComponent, State, Prop } from '@bearer/core'
import '@bearer/ui'

@RootComponent({
  group: 'setup',
  name: 'action'
})
export class SetupAction {
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

