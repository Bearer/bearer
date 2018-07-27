/*
  The purpose of this component is to save scenario credentials.
*/

import { Component, Prop } from "@bearer/core";
import "@bearer/ui";

@Component({
  tag: "{{componentTagName}}-setup-display",
  shadow: true
})
export class {{componentName}}SetupDisplay {
  render() {
    return (
      <bearer-setup-display scenarioId="BEARER_SCENARIO_ID" />
    )
  }
  @Prop() BEARER_ID: string;
}

