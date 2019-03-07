/*
  The purpose of this component is to save integration credentials.
*/

import { RootComponent } from "@bearer/core";
import "@bearer/ui";

@RootComponent({
  name: 'setup-view',
})
export class SetupView {
  render() {
    return (
      <bearer-setup-display integrationId="BEARER_INTEGRATION_ID" setupId={(this as any).setupId} />
    )
  }
}
