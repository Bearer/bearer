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
      if (!needProcessing(tsSourceFile)) {
        return tsSourceFile
      }

      // Inject Imports if needed: Watch
      const preparedSourceFile = ensureInjectedWatchDecorator(tsSourceFile)

      function visit(node: ts.Node): ts.VisitResult<ts.Node> {
        if (ts.isClassDeclaration(node)) {
          // Ensures we have context available
          const withInjectedContext = ensureInjectedContext(node)

          // Inject prop watcher
          const withInjectedWatcher = injectPropertyWatcher(withInjectedContext)

          // Append logic to componentWillLoad/componentDidUnload
          const withComponentLifecyleHooked = updateComponentLifecycle(
            withInjectedWatcher
          )
          // Add update logic method
          const bearerStateReadyComponent = injectStateUpdateLogic(
            withComponentLifecyleHooked
          )
          return bearerStateReadyComponent
        }
        return ts.visitEachChild(node, visit, transformContext)
      }

      return visit(preparedSourceFile) as ts.SourceFile
    }
  }
}

function injectStateUpdateLogic(
  classNode: ts.ClassDeclaration
): ts.ClassDeclaration {
  // updateFromState = state => {
  //   this.attachedPullRequests = state['attachedPullRequests']
  // }
  return classNode
}

function ensureInjectedContext(
  classNode: ts.ClassDeclaration
): ts.ClassDeclaration {
  return classNode
}

function updateComponentLifecycle(
  classNode: ts.ClassDeclaration
): ts.ClassDeclaration {
  // componentWillLoad() {
  //   this.context.subscribe(this)
  // }
  // componentDidUnload() {
  //   this.context.unsubscribe(this)
  // }
  return classNode
}

function injectPropertyWatcher(
  classNode: ts.ClassDeclaration
): ts.ClassDeclaration {
  // overrides if one exist
  // @Watch('attachedPullRequests')
  // changeRepo(newValue: any) {
  //   console.log('[BEARER]', 'attachedPullRequests updated')
  //   this.context.update('attachedPullRequests', newValue)
  // }
  return classNode
}

function ensureInjectedWatchDecorator(
  sourceFile: ts.SourceFile
): ts.SourceFile {
  return sourceFile
}

/**
 *  Not a declaration file and contains a @BearerState propertyDecorator
 */
function needProcessing(sourceFile: ts.SourceFile): boolean {
  if (sourceFile.isDeclarationFile) {
    return false
  }

  const shouldProcess = false

  return shouldProcess
}
