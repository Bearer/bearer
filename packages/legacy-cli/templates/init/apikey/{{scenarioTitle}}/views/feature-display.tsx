/*
  The purpose of this component is to be the result of your scenario.
  Its responsibility is to retrieve the scenario state from a previous action
  of a user.
*/
import { RootComponent } from '@bearer/core'
import '@bearer/ui'

@RootComponent({
  name: 'display',
  group: 'feature'
})
export class FeatureDisplay {
  render() {
    return <bearer-alert kind="success">🚀 My {{scenarioTitle}} Display component 🚀</bearer-alert>
  }
}
