/*
 *
 */
import * as ts from 'typescript'
import { propDecoratedWithName } from './decorator-helpers'
import { ensureWatchImported } from './bearer'
import { Decorators, Component } from './constants'

type TransformerOptions = {
  verbose?: true
}

const state = ts.createIdentifier('state')

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
 * Inject methods to auto update component
 */
function injectStateUpdateLogic(
  classNode: ts.ClassDeclaration
): ts.ClassDeclaration {
  // TODO : dynamic propName / dynamic statePropName
  const propName = 'attachedPullRequests'
  const statePropName = 'attachedPullRequests'
  return ts.updateClassDeclaration(
    classNode,
    classNode.decorators,
    classNode.modifiers,
    classNode.name,
    classNode.typeParameters,
    classNode.heritageClauses,
    [
      ...classNode.members,
      ts.createProperty(
        undefined,
        undefined,
        'bearerUpdateFromState',
        undefined,
        undefined,
        ts.createArrowFunction(
          undefined,
          undefined,
          [ts.createParameter(undefined, undefined, undefined, state)],
          undefined,
          undefined,
          ts.createBlock([
            ts.createStatement(
              ts.createAssignment(
                ts.createPropertyAccess(ts.createThis(), propName),
                ts.createElementAccess(state, ts.createLiteral(statePropName))
              )
            )
          ])
        )
      )
    ]
  )
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
    ts.createPropertyAccess(
      ts.createThis(),
      `${Component.COMPONENT_BEARER_CONTEXT_PROP}.subscribe`
    ),
    undefined,
    [ts.createThis()]
  )

  const componentDidUnload = ts.createCall(
    ts.createPropertyAccess(
      ts.createThis(),
      `${Component.COMPONENT_BEARER_CONTEXT_PROP}.unsubscribe`
    ),
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
        Component.componentWillLoad,
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
        Component.componentDidUnload,
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
              ts.createIdentifier(Decorators.Watch) as ts.Expression,
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
              ts.createPropertyAccess(
                ts.createThis(),
                `${Component.COMPONENT_BEARER_CONTEXT_PROP}.update`
              ),
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
        Decorators.BearerState
      )
    }
  })

  return shouldProcess
}
