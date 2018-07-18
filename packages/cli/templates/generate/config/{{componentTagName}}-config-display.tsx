/*
  The purpose of this component is to save scenario configuration.
*/

import { Component, Prop } from "@bearer/core";
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
  @Prop() BEARER_ID: string;
}

