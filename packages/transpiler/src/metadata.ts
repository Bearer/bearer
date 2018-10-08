import * as ts from 'typescript'

import { ComponentMetadata } from './types'

function relativePath(path: string): string {
  return path.replace(process.cwd(), '')
}
export default class Metadata {
  components: Array<ComponentMetadata> = []

  constructor(readonly prefix: string = null, readonly suffix: string = null) {}

  registerComponent = (component: ComponentMetadata): void => {
    this.components = [
      ...this.components.filter(otherComponent => component.finalTagName !== otherComponent.finalTagName),
      {
        ...component,
        fileName: relativePath(component.fileName)
      }
    ].sort((a, b) => (a.finalTagName > b.finalTagName ? 1 : -1))
  }

  findComponentFrom(tsSourceFile: ts.SourceFile): ComponentMetadata {
    return this.components.find(component => component.fileName === relativePath(tsSourceFile.fileName))
  }
}
