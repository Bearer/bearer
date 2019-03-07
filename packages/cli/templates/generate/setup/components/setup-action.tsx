/*
  The purpose of this component is to save integration credentials.
  This file has been generated automatically. Edit it at your own risk :-)
*/

import { RootComponent, State, Prop, Output } from '@bearer/core'
import '@bearer/ui'
import { FieldSet } from '@bearer/ui/lib/collection/components/Forms/Fieldset'

@RootComponent({
  name: 'setup-action',
})
export class SetupAction {
  @Prop() display: 'inline' | 'block' = 'inline'

  @State() fields = new FieldSet({{fields}})

  @Output() setup: any

  setupSubmitHandler = (e: CustomEvent) => {
    this.setup = e.detail
  }

  render() {
    return (
      <bearer-setup
        display={this.display}
        onSetupSubmit={this.setupSubmitHandler}
        integrationId='BEARER_INTEGRATION_ID'
        fields={this.fields}
      />
    )
  }
}
