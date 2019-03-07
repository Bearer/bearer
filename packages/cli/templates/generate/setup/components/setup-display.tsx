/*
  The purpose of this component is to save scenario credentials.
*/

import { RootComponent } from "@bearer/core";
import "@bearer/ui";

@RootComponent({
  name: 'setup-display',
})
export class SetupDisplay {
  render() {
    return (
      <bearer-setup-display scenarioId="BEARER_SCENARIO_ID" setupId={(this as any).setupId} />
    )
  }
}
