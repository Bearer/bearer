/*
  The purpose of this component is to deal with scenario navigation between each screens.

*/

import { Component } from '@bearer/core'
import '@bearer/ui'

@Component({
  tag: '{{componentTagName}}',
  styleUrl: '{{scenarioTitle}}.css',
  shadow: true
})
export class {{scenarioTitle}}Action {
  render() {
    return (
      <div>
        <bearer-navigator>
        </bearer-navigator>
      </div>
    )
  }
}
