/*
 *
 */
import * as ts from 'typescript'

import { Component, Decorators } from '../constants'
import { hasDecoratorNamed, hasPropDecoratedWithName, propDecoratedWithName } from '../helpers/decorator-helpers'
import { ensureMethodExists } from '../helpers/guards-helpers'
import { prependToStatements, updateMethodOfClass } from '../helpers/method-updaters'
import { isWatcherOn } from '../helpers/stencil-helpers'
import { TransformerOptions } from '../types'

import { ensureBearerContextInjected, ensureStateImported, ensureWatchImported } from './bearer'

const state = ts.createIdentifier('state')

export default function BearerStateInjector({  }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return transformContext => {
    return tsSourceFile => {
      if (!hasBearerStateDecorator(tsSourceFile)) {
        return tsSourceFile
      }

      const propsDecorator = extractDecoratedPropertyInformation(tsSourceFile)
      // Inject Imports if needed: Watch
      const preparedSourceFile = ensureStateImported(ensureWatchImported(tsSourceFile))

      function visit(node: ts.Node): ts.VisitResult<ts.Node> {
        if (ts.isClassDeclaration(node)) {
          // Ensures we have context available
          const withInjectedContext = ensureBearerContextInjected(node)

          // Append @Prop() decorator to before @BearerState
          const withPropDecoratorToDeclaration = addPropDecoratorToPropDeclaration(withInjectedContext)

          // Inject prop watcher
          const withInjectedWatcher = injectPropertyWatcher(withPropDecoratorToDeclaration, propsDecorator)

          // Append logic to componentWillLoad/componentDidUnload
          const withComponentLifecyleHooked = updateComponentLifecycle(withInjectedWatcher)
          // Add update logic method
          const bearerStateReadyComponent = injectStateUpdateLogic(withComponentLifecyleHooked, propsDecorator)

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
                  ts.createPropertyAccess(ts.createThis(), meta.componentPropName),
                  ts.createElementAccess(state, ts.createLiteral(meta.statePropName))
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
function updateComponentLifecycle(aClassNode: ts.ClassDeclaration): ts.ClassDeclaration {
  const classNode = ensureMethodExists(
    ensureMethodExists(aClassNode, Component.componentWillLoad),
    Component.componentDidUnload
  )

  const withSubscribe = updateMethodOfClass(classNode, Component.componentWillLoad, method =>
    prependToStatements(method, [
      ts.createStatement(
        ts.createCall(ts.createPropertyAccess(ts.createThis(), `${Component.bearerContext}.subscribe`), undefined, [
          ts.createThis()
        ])
      )
    ])
  )

  return updateMethodOfClass(withSubscribe, Component.componentDidUnload, method =>
    prependToStatements(method, [
      ts.createStatement(
        ts.createCall(ts.createPropertyAccess(ts.createThis(), `${Component.bearerContext}.unsubscribe`), undefined, [
          ts.createThis()
        ])
      )
    ])
  )
}

function createWatcher(meta: IDecoratedPropInformation): ts.MethodDeclaration {
  return ts.createMethod(
    [
      ts.createDecorator(
        ts.createCall(ts.createIdentifier(Decorators.Watch) as ts.Expression, undefined, [
          ts.createLiteral(meta.componentPropName)
        ])
      )
    ],
    undefined,
    undefined,
    ts.createIdentifier(`_notifyBearerStateHandler_${meta.componentPropName}`),
    undefined,
    undefined,
    [
      ts.createParameter(
        undefined,
        undefined,
        undefined,
        'newValue',
        undefined,
        ts.createKeywordTypeNode(ts.SyntaxKind.AnyKeyword),
        undefined
      )
    ],
    undefined,
    ts.createBlock([createUpdateStatement(meta.statePropName, 'newValue')])
  )
}

function createUpdateStatement(stateName: string, parameterName: string): ts.Statement {
  return ts.createStatement(
    ts.createCall(ts.createPropertyAccess(ts.createThis(), `${Component.bearerContext}.update`), undefined, [
      ts.createLiteral(stateName),
      ts.createIdentifier(parameterName)
    ])
  )
}
/**
 * Add or update State Watcher
 */
function injectPropertyWatcher(
  classNode: ts.ClassDeclaration,
  propsDecoratedMeta: Array<IDecoratedPropInformation>
): ts.ClassDeclaration {
  const members = propsDecoratedMeta.reduce((members, meta) => {
    const predicate = node => ts.isMethodDeclaration(node) && isWatcherOn(node, meta.componentPropName)
    const watcherHandler = members.find(predicate) as ts.MethodDeclaration
    if (watcherHandler) {
      const newValueParameterName = (watcherHandler.parameters[0].name as ts.Identifier).escapedText as string
      return [
        ...members.filter(node => !predicate(node)),
        prependToStatements(watcherHandler, [createUpdateStatement(meta.componentPropName, newValueParameterName)])
      ]
    }
    return ts.createNodeArray([...members, createWatcher(meta)])
  }, classNode.members)
  return ts.updateClassDeclaration(
    classNode,
    classNode.decorators,
    classNode.modifiers,
    classNode.name,
    classNode.typeParameters,
    classNode.heritageClauses,
    members
  )
}
/**
 * Add @Prop() before @BearerState
 * withPropDecoratorToDeclaration
 */

function addPropDecoratorToPropDeclaration(classNode: ts.ClassDeclaration): ts.ClassDeclaration {
  return ts.updateClassDeclaration(
    classNode,
    classNode.decorators,
    classNode.modifiers,
    classNode.name,
    classNode.typeParameters,
    classNode.heritageClauses,
    classNode.members.map(appendStateDecoratorIfNeeded)
  )
}

function appendStateDecoratorIfNeeded(element: ts.ClassElement): ts.ClassElement {
  if (
    ts.isPropertyDeclaration(element) &&
    hasDecoratorNamed(element as ts.PropertyDeclaration, Decorators.BearerState)
  ) {
    return ts.updateProperty(
      element,
      [
        ...element.decorators,
        ts.createDecorator(ts.createCall(ts.createIdentifier(Decorators.State), undefined, undefined))
      ],
      element.modifiers,
      element.name,
      element.questionToken,
      element.type,
      element.initializer
    )
  }
  return element
}

/**
 *  Not a declaration file and contains a @BearerState propertyDecorator
 */
function hasBearerStateDecorator(sourceFile: ts.SourceFile): boolean {
  if (sourceFile.isDeclarationFile) {
    return false
  }

  return ts.forEachChild(
    sourceFile,
    node => ts.isClassDeclaration(node) && hasPropDecoratedWithName(node, Decorators.BearerState)
  )
}

type IDecoratedPropInformation = {
  componentPropName: string
  statePropName: string
}

function extractDecoratedPropertyInformation(tsSourceFile: ts.SourceFile): Array<IDecoratedPropInformation> {
  return (
    ts.forEachChild(
      tsSourceFile,
      node => ts.isClassDeclaration(node) && propDecoratedWithName(node, Decorators.BearerState)
    ) || []
  ).map(prop => {
    const decoratorOptions = (prop.decorators[0].expression as ts.CallExpression).arguments[0]
    const componentPropName: string = (prop.name as ts.Identifier).escapedText as string

    let statePropName = componentPropName

    if (decoratorOptions && ts.isObjectLiteralExpression(decoratorOptions)) {
      const stateNameOption = decoratorOptions.properties.find(
        prop => (prop.name as ts.Identifier).escapedText === Decorators.statePropName
      ) as ts.PropertyAssignment
      if (stateNameOption) {
        statePropName = (stateNameOption.initializer as ts.StringLiteral).text
      }
    }
    return {
      componentPropName,
      statePropName
    }
  })
}
