/*
  The purpose of this component is to be the result of your scenario.
  Its responsibility is to retrieve the scenario state from a previous action
  of a user.
*/
import { Component } from '@bearer/core'
import '@bearer/ui'

@Component({
  tag: '{{componentTagName}}-display',
  styleUrl: '{{scenarioTitle}}.css',
  shadow: true
})
export class {{scenarioTitle}}Display {
  render() {
    return (
      <bearer-typography as="h1" kind="h1">🎉🎉 display component 🎉🎉</bearer-typography>
    )
  }
}
