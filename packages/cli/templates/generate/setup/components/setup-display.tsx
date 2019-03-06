/*
  The purpose of this component is to save integration credentials.
*/

import { RootComponent } from "@bearer/core";
import "@bearer/ui";

@RootComponent({
  group: "setup",
  role: 'display'
})
export class SetupDisplay {
  render() {
    return (
      <bearer-setup-display integrationId="BEARER_INTEGRATION_ID" setupId={(this as any).setupId} />
    )
  }
}
