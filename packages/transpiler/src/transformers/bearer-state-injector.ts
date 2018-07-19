/*
 *
 */
import * as ts from 'typescript'
import { propDecoratedWithName } from './decorator-helpers'
import { ensureWatchImported } from './bearer'

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
      const preparedSourceFile = ensureWatchImported(tsSourceFile)

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

/**
 * TODO
 *
 */
function injectStateUpdateLogic(
  classNode: ts.ClassDeclaration
): ts.ClassDeclaration {
  // updateFromState = state => {
  //   this.attachedPullRequests = state['attachedPullRequests']
  // }
  return classNode
}

/**
 * TODO
 */

function ensureInjectedContext(
  classNode: ts.ClassDeclaration
): ts.ClassDeclaration {
  return classNode
}

/**
 * Add subscription methods to component lifecycle
 */
function updateComponentLifecycle(
  classNode: ts.ClassDeclaration
): ts.ClassDeclaration {
  // TODO: override and append if it exists
  const componentWillLoad = ts.createCall(
    ts.createPropertyAccess(ts.createThis(), 'context.subscribe'),
    undefined,
    [ts.createThis()]
  )

  const componentDidUnload = ts.createCall(
    ts.createPropertyAccess(ts.createThis(), 'context.unsubscribe'),
    undefined,
    [ts.createThis()]
  )

  return ts.updateClassDeclaration(
    classNode,
    classNode.decorators,
    classNode.modifiers,
    classNode.name,
    classNode.typeParameters,
    classNode.heritageClauses,
    [
      ...classNode.members,
      ts.createMethod(
        undefined,
        undefined,
        undefined,
        'componentWillLoad',
        undefined,
        undefined,
        undefined,
        undefined,
        ts.createBlock([ts.createStatement(componentWillLoad)])
      ),
      ts.createMethod(
        undefined,
        undefined,
        undefined,
        'componentDidUnload',
        undefined,
        undefined,
        undefined,
        undefined,
        ts.createBlock([ts.createStatement(componentDidUnload)])
      )
    ]
  )
}

/**
 * Add or update State Watcher
 */
function injectPropertyWatcher(
  classNode: ts.ClassDeclaration
): ts.ClassDeclaration {
  // TODO: make it dynamic
  const propName = 'propName'
  // TODO: override if one already exist
  return ts.updateClassDeclaration(
    classNode,
    classNode.decorators,
    classNode.modifiers,
    classNode.name,
    classNode.typeParameters,
    classNode.heritageClauses,
    [
      ...classNode.members,
      ts.createMethod(
        [
          ts.createDecorator(
            ts.createCall(
              ts.createIdentifier('Watch') as ts.Expression,
              undefined,
              [ts.createLiteral(propName)]
            )
          )
        ],
        undefined,
        undefined,
        ts.createIdentifier('_notifyBearerStateHandler'),
        undefined,
        undefined,
        undefined,
        undefined,
        ts.createBlock([
          ts.createStatement(
            ts.createCall(
              ts.createPropertyAccess(ts.createThis(), 'context.update'),
              undefined,
              [ts.createLiteral(propName), ts.createThis()]
            )
          )
        ])
      )
    ]
  )
}

/**
 *  Not a declaration file and contains a @BearerState propertyDecorator
 */
function needProcessing(sourceFile: ts.SourceFile): boolean {
  if (sourceFile.isDeclarationFile) {
    return false
  }

  let shouldProcess = false
  ts.forEachChild(sourceFile, node => {
    if (ts.isClassDeclaration(node)) {
      shouldProcess = propDecoratedWithName(
        node as ts.ClassDeclaration,
        'BearerState'
      )
    }
  })

  return shouldProcess
}
