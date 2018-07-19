/*
 *
 */
import * as ts from 'typescript'
import { propDecoratedWithName } from './decorator-helpers'
type TransformerOptions = {
  verbose?: true
}

export default function BearerStateInjector({
  verbose
}: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return transformContext => {
    return tsSourceFile => {
      function visit(node: ts.Node): ts.VisitResult<ts.Node> {
        if (
          ts.isClassDeclaration(node) &&
          propDecoratedWithName(node as ts.ClassDeclaration, 'BearerState')
        ) {
          // check if context already present
          //  ensure context injected
          //
          //  override componentDidLoad
          //  override componentDidLoad
          // implement updateFromState
          //
          //  inject : import watch
          //  implement watcher
          // @Watch('attachedPullRequests')
          // changeRepo(newValue: any) {
          //   console.log('[BEARER]', 'attachedPullRequests updated')
          //   this.context.update('attachedPullRequests', newValue)
          // }
          // updateFromState = state => {
          //   this.attachedPullRequests = state['attachedPullRequests']
          // }
          // componentWillLoad() {
          //   this.context.subscribe(this)
          // }
          // componentDidUnload() {
          //   this.context.unsubscribe(this)
          // }
        }
        return ts.visitEachChild(node, visit, transformContext)
      }

      return visit(tsSourceFile) as ts.SourceFile
    }
  }
}
