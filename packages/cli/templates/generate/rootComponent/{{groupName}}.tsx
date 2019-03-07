/*
  The purpose of this component is to deal with integration navigation between each views.

*/

import { RootComponent } from '@bearer/core'
import '@bearer/ui'

@RootComponent({
  name: '{{groupName}}',
})
export class {{componentClassName}} {
  render() {
    return (
      <bearer-navigator direction="right">
        <span slot="navigator-btn-content">{{componentName}} Action</span>
        {{withAuthScreen}}
        <bearer-navigator-screen navigationTitle="My first screen">
          <div style={ { textAlign: 'center' } }>My first screen</div>
        </bearer-navigator-screen>
      </bearer-navigator>
    )
  }
}
