/*
  The purpose of this component is to be the result of your integration.
  Its responsibility is to retrieve the integration state from a previous action
  of a user.
*/
import { RootComponent } from '@bearer/core'
import '@bearer/ui'

@RootComponent({
  role: 'display',
  group: '{{groupName}}'
})
export class {{componentClassName}}Display {
  render() {
    return <bearer-alert kind="success">🚀 My {{componentName}} Display component 🚀</bearer-alert>
  }
}
