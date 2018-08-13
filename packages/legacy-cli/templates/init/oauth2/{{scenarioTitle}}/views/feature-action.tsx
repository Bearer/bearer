/*
  The purpose of this component is to deal with scenario navigation between each views.

*/

import { RootComponent } from '@bearer/core'
import '@bearer/ui'

@RootComponent({
  name: 'action',
  group: 'feature'
})
export class FeatureAction {
  render() {
    return (
      <bearer-navigator btnProps={ { content:"{{scenarioTitle}}", kind:"primary" } } direction="right">
        <bearer-navigator-auth-screen />
        <bearer-navigator-screen navigationTitle="My first screen">
          <div style={ { textAlign: 'center' } }>My first screen</div>
        </bearer-navigator-screen>
      </bearer-navigator>
    )
  }
}
