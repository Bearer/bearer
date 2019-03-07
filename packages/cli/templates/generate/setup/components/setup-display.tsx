/*
  The purpose of this component is to display the integration reference ID.
  This file has been generated automatically. Edit it at your own risk :-)
*/

import { RootComponent, Input } from '@bearer/core'
import '@bearer/ui'

@RootComponent({
  group: 'setup',
  role: 'display'
})
export class SetupDisplay {
  @Input() setup: any

  render() {
    return <bearer-setup-display integrationId="BEARER_INTEGRATION_ID" setupId={this.setup && this.setup.referenceId} />
  }
}
