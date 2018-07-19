/*
 *
 */
import * as ts from 'typescript'
import {
  propDecoratedWithName,
  hasPropDecoratedWithName
} from './decorator-helpers'
import { ensureWatchImported, ensureBearerContextInjected } from './bearer'
import { Decorators, Component } from './constants'

/**
 * TODOS:
 *  * add typing on newValue parameter of the watch handler
 */

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

      const propsDecorator = extractDecoratedPropertyInformation(tsSourceFile)
      // Inject Imports if needed: Watch
      const preparedSourceFile = ensureWatchImported(tsSourceFile)

      function visit(node: ts.Node): ts.VisitResult<ts.Node> {
        if (ts.isClassDeclaration(node)) {
          // Ensures we have context available
          const withInjectedContext = ensureBearerContextInjected(node)

          // Inject prop watcher
          const withInjectedWatcher = injectPropertyWatcher(
            withInjectedContext,
            propsDecorator
          )

          // Append logic to componentWillLoad/componentDidUnload
          const withComponentLifecyleHooked = updateComponentLifecycle(
            withInjectedWatcher
          )
          // Add update logic method
          const bearerStateReadyComponent = injectStateUpdateLogic(
            withComponentLifecyleHooked,
            propsDecorator
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
  classNode: ts.ClassDeclaration,
  propsDecoratedMeta: Array<IDecoratedPropInformation>
): ts.ClassDeclaration {
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
          ts.createBlock(
            propsDecoratedMeta.map(meta =>
              ts.createStatement(
                ts.createAssignment(
                  ts.createPropertyAccess(
                    ts.createThis(),
                    meta.componentPropName
                  ),
                  ts.createElementAccess(
                    state,
                    ts.createLiteral(meta.statePropName)
                  )
                )
              )
            )
          )
        )
      )
    ]
  )
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
      `${Component.bearerContext}.subscribe`
    ),
    undefined,
    [ts.createThis()]
  )

  const componentDidUnload = ts.createCall(
    ts.createPropertyAccess(
      ts.createThis(),
      `${Component.bearerContext}.unsubscribe`
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
  classNode: ts.ClassDeclaration,
  propsDecoratedMeta: Array<IDecoratedPropInformation>
): ts.ClassDeclaration {
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
      ...propsDecoratedMeta.map(meta =>
        ts.createMethod(
          [
            ts.createDecorator(
              ts.createCall(
                ts.createIdentifier(Decorators.Watch) as ts.Expression,
                undefined,
                [ts.createLiteral(meta.componentPropName)]
              )
            )
          ],
          undefined,
          undefined,
          ts.createIdentifier('_notifyBearerStateHandler'),
          undefined,
          undefined,
          [
            ts.createParameter(
              undefined,
              undefined,
              undefined,
              'newValue',
              undefined,
              undefined,
              undefined
            )
          ],
          undefined,
          ts.createBlock([
            ts.createStatement(
              ts.createCall(
                ts.createPropertyAccess(
                  ts.createThis(),
                  `${Component.bearerContext}.update`
                ),
                undefined,
                [
                  ts.createLiteral(meta.statePropName),
                  ts.createIdentifier('newValue')
                ]
              )
            )
          ])
        )
      )
    ]
  )
}

// function getDecoratedProp(classNode: ts.ClassDeclaration): null |
/**
 *  Not a declaration file and contains a @BearerState propertyDecorator
 */
function needProcessing(sourceFile: ts.SourceFile): boolean {
  if (sourceFile.isDeclarationFile) {
    return false
  }

  return ts.forEachChild(
    sourceFile,
    node =>
      ts.isClassDeclaration(node) &&
      hasPropDecoratedWithName(node, Decorators.BearerState)
  )
}
interface IDecoratedPropInformation {
  componentPropName: string
  statePropName: string
}
function extractDecoratedPropertyInformation(
  tsSourceFile: ts.SourceFile
): Array<IDecoratedPropInformation> {
  return (
    ts.forEachChild(
      tsSourceFile,
      node =>
        ts.isClassDeclaration(node) &&
        propDecoratedWithName(node, Decorators.BearerState)
    ) || []
  ).map(prop => {
    const decoratorOptions = (prop.decorators[0]
      .expression as ts.CallExpression).arguments[0]
    const componentPropName: string = prop.name['escapedText']

    let statePropName = componentPropName

    if (decoratorOptions && ts.isObjectLiteralExpression(decoratorOptions)) {
      const stateNameOption = decoratorOptions.properties.find(
        prop => prop.name['escapedText'] == 'statePropName'
      )
      if (stateNameOption) {
        statePropName = stateNameOption['initializer']
      }
    }
    return {
      componentPropName,
      statePropName
    }
  })
}
