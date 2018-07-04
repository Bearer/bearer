/*
  The purpose of this component is to save scenario credentials.
*/

import { Component } from "@bearer/core";
import "@bearer/ui";

@Component({
  tag: "{{componentTagName}}-setup-display",
  shadow: true
})
export class {{scenarioTitle}}SetupDisplay {
  render() {
    return (
      <bearer-setup-display scenarioId="BEARER_SCENARIO_ID" />
    )
  }
}

