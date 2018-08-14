/*
  The purpose of this component is to be the result of your scenario.
  Its responsibility is to retrieve the scenario state from a previous action
  of a user.
*/
import { RootComponent } from '@bearer/core'
import '@bearer/ui'

@RootComponent({
  name: 'display',
  group: '{{groupName}}'
})
export class {{componentClassName}}Display {
  render() {
    return <bearer-alert kind="success">🚀 My {{componentName}} Display component 🚀</bearer-alert>
  }
}
