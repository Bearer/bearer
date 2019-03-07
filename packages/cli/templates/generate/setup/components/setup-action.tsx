/*
  The purpose of this component is to save integration credentials.
  This file has been generated automatically and should not be edited.
*/

import Bearer, { RootComponent, State, Prop, Output, Element, Listen } from '@bearer/core'
import '@bearer/ui'
import { FieldSet } from "@bearer/ui/lib/collection/components/Forms/Fieldset";

@RootComponent({
  name: 'setup-action',
})
export class SetupAction {
  @Prop() onSetupSuccess: (detail: any) => void = (_any: any) => { }
  @State() fields = new FieldSet({{fields}})
  @State() innerListener = `setup_success:BEARER_INTEGRATION_ID`
  @Output() setup: any;
  @Element() el: HTMLElement;

  @Listen("setup-setupSaved")
  setupSavedHandler(e: CustomEvent) {
    const event = new CustomEvent("setupSuccess", e);
    document.dispatchEvent(event);
    this.el.shadowRoot
      .querySelector<HTMLBearerDropdownButtonElement>("bearer-dropdown-button")
      .toggle(false);
    Bearer.emitter.emit(this.innerListener, {
      referenceId: e.detail.referenceId
    });
  }

  handleSubmit = (e: CustomEvent) => {
    this.setup = e.detail.set.reduce((acc, obj) => ({ ...acc, [obj.controlName]: obj.value }), {ReadAllowed: false})
  };

  render() {
    return (
      <bearer-dropdown-button>
        <span slot="dropdown-btn-content">Setup component</span>
        <bearer-form fields={this.fields} onSubmit={this.handleSubmit} />
      </bearer-dropdown-button>
    )
  }
}
