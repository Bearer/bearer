import { ComponentMetadata } from './types'

export default class Metadata {
  components: Array<ComponentMetadata> = []

  constructor(readonly prefix: string = null, readonly suffix: string = null) {}

  registerComponent = (component: ComponentMetadata): void => {
    this.components = [
      ...this.components.filter(otherComponent => component.finalTagName !== otherComponent.finalTagName),
      component
    ].sort((a, b) => (a.finalTagName > b.finalTagName ? 1 : -1))
  }
}
