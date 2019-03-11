/*
  The purpose of this component is to show integration setup id.
  This file has been generated automatically. Edit it at your own risk :-)
*/

import { RootComponent, Input } from "@bearer/core";
import "@bearer/ui";

@RootComponent({
  name: 'setup-view',
})
export class SetupView {
  @Input() setup: any

  render() {
    return (
      <bearer-setup-display integrationId="BEARER_INTEGRATION_ID" setupId={(this as any).setupId} />
    )
  }
}
