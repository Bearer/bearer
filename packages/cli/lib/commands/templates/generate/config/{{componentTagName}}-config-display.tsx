/*
  The purpose of this component is to save scenario configuration.
*/

import { Component } from "@bearer/core";
import "@bearer/ui";

@Component({
  tag: "{{componentTagName}}-config-display",
  shadow: true
})
export class {{scenarioTitle}}ConfigDisplay {
  render() {
    return (
      <bearer-config-display scenarioId="BEARER_SCENARIO_ID" />
    )
  }
}

