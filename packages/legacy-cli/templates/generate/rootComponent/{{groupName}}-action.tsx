/*
  The purpose of this component is to deal with scenario navigation between each views.

*/

import { RootComponent } from '@bearer/core'
import '@bearer/ui'

@RootComponent({
  name: 'action',
  group: '{{groupName}}'
})
export class {{groupName}}Action {
  render() {
    return (
      <div>
        <bearer-navigator btnProps={ { content:"{{componentName}} Action", kind:"primary" } } direction="right">
          <bearer-navigator-auth-screen />
        </bearer-navigator>
      </div>
    )
  }
}
