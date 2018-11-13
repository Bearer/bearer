/*
 * Input Transformer
 */
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { hasDecoratorNamed } from '../helpers/decorator-helpers'
import { getNodeName } from '../helpers/node-helpers'
import { TransformerOptions } from '../types'

import {
  ensureIntentImported,
  ensureListenImported,
  ensurePropImported,
  ensureStateImported,
  ensureWatchImported
} from './bearer'

export default function InputDecorator(_options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    return tsSourceFile => {
      if (tsSourceFile.isDeclarationFile) {
        return tsSourceFile
      }
      // gather inputs information
      const inputsMeta = retrieveInputsMetas(tsSourceFile)

      if (!inputsMeta.length) {
        return tsSourceFile
      }
      const sourceFileWithImports = [
        ensureListenImported,
        ensureStateImported,
        ensureIntentImported,
        ensureWatchImported,
        ensurePropImported
      ].reduce((sourceFile, importer) => importer(sourceFile), tsSourceFile)

      return ts.visitEachChild(sourceFileWithImports, replaceInputVisitor(inputsMeta), _transformContext)
    }

    function retrieveInputsMetas(sourcefile: ts.SourceFile): Array<InputMeta> {
      const inputs: Array<InputMeta> = []

      const visitor = (tsNode: ts.Node) => {
        if (ts.isPropertyDeclaration(tsNode) && hasDecoratorNamed(tsNode, Decorators.Input)) {
          const name = getNodeName(tsNode)
          const capitalizedName = name.charAt(0).toUpperCase() + name.slice(1)
          inputs.push({
            propDeclarationName: name,
            scope: 'string', // TODO: retrieve from options
            propName: `${name}RefId`, // TODO: retrieve from options
            eventName: `${name}:saved`, // TODO: retrieve from options
            intentName: `get${capitalizedName}`, // TODO: retrieve from options
            intentMethodName: `fetcherGet${capitalizedName}`, // TODO: retrieve from options
            autoUpdate: true, // TODO: retrieve from options
            loadMethodName: `_load${capitalizedName}`,
            typeIdentifier: tsNode.type,
            intializer: tsNode.initializer,
            watcherName: `_watch${capitalizedName}`
          })
        }
        return ts.visitEachChild(tsNode, visitor, _transformContext)
      }

      ts.visitEachChild(sourcefile, visitor, _transformContext)

      return inputs
    }

    function injectInputStatements(tsClass: ts.ClassDeclaration, inputsMeta: Array<InputMeta>): ts.ClassDeclaration {
      const newMembers = inputsMeta.reduce(
        (members, meta) => {
          // create @State()
          const inputMembers = [
            createLocalStateProperty(meta),
            createRefIdProp(meta),
            createEventListener(meta),
            createLoadResourceMethod(meta),
            createFetcher(meta),
            createRefIdWatcher(meta)
          ]
          return members.concat(inputMembers)
        },
        [...tsClass.members]
      )

      return ts.updateClassDeclaration(
        tsClass,
        tsClass.decorators,
        tsClass.modifiers,
        tsClass.name,
        tsClass.typeParameters,
        tsClass.heritageClauses,
        newMembers
      )
    }

    function replaceInputVisitor(inputsMeta: Array<InputMeta>): (tsNode: ts.Node) => ts.VisitResult<ts.Node> {
      return (tsNode: ts.Node) => {
        // remove input usage
        if (ts.isPropertyDeclaration(tsNode) && hasDecoratorNamed(tsNode, Decorators.Input)) {
          return null
        }

        if (ts.isClassDeclaration(tsNode)) {
          return ts.visitEachChild(
            injectInputStatements(tsNode, inputsMeta),
            replaceInputVisitor(inputsMeta),
            _transformContext
          )
        }
        return ts.visitEachChild(tsNode, replaceInputVisitor(inputsMeta), _transformContext)
      }
    }
  }
}

function createLocalStateProperty(meta: InputMeta) {
  return ts.createProperty(
    [ts.createDecorator(ts.createCall(ts.createIdentifier(Decorators.State), undefined, undefined))],
    undefined,
    ts.createIdentifier(meta.propDeclarationName),
    undefined,
    meta.typeIdentifier,
    meta.intializer
  )
}

function createRefIdProp(meta: InputMeta) {
  return ts.createProperty(
    [
      ts.createDecorator(
        ts.createCall(ts.createIdentifier(Decorators.Prop), undefined, [
          ts.createObjectLiteral([ts.createPropertyAssignment(ts.createLiteral('mutable'), ts.createTrue())])
        ])
      )
    ],
    undefined,
    ts.createIdentifier(meta.propName),
    undefined,
    ts.createKeywordTypeNode(ts.SyntaxKind.StringKeyword),
    undefined
  )
}

// This will make the prop change when an event is triggered
function createEventListener(meta: InputMeta) {
  return ts.createMethod(
    [
      ts.createDecorator(
        ts.createCall(ts.createIdentifier(Decorators.Listen), undefined, [ts.createLiteral(meta.eventName)])
      )
    ],
    undefined,
    undefined,
    `${meta.propName}Changed`,
    undefined,
    undefined,
    [ts.createParameter(undefined, undefined, undefined, ts.createIdentifier('event'), undefined, undefined)],
    undefined,
    ts.createBlock([
      ts.createStatement(
        ts.createBinary(
          ts.createPropertyAccess(ts.createThis(), meta.propName),
          ts.SyntaxKind.EqualsToken,
          ts.createIdentifier('event.detail.referenceId')
        )
      )
    ])
  )
}

function createLoadResourceMethod(meta: InputMeta) {
  const intentCall = ts.createCall(
    ts.createPropertyAccess(ts.createThis(), meta.intentMethodName),
    undefined,
    undefined
  )
  const udapteState = ts.createArrowFunction(
    undefined,
    undefined,
    [
      ts.createParameter(
        undefined,
        undefined,
        undefined,
        ts.createObjectBindingPattern([ts.createBindingElement(undefined, undefined, 'data')]),
        undefined,
        ts.createTypeLiteralNode([
          ts.createPropertySignature(undefined, 'data', undefined, meta.typeIdentifier, undefined)
        ]),
        undefined
      )
    ],
    undefined,
    undefined,
    ts.createBlock([
      ts.createStatement(
        ts.createBinary(
          ts.createPropertyAccess(ts.createThis(), meta.propDeclarationName),
          ts.SyntaxKind.EqualsToken,
          ts.createIdentifier('data')
        )
      )
    ])
  )
  const promiseHandler = ts.createCall(ts.createPropertyAccess(intentCall, 'then'), undefined, [udapteState])
  return ts.createProperty(
    undefined,
    undefined,
    meta.loadMethodName,
    undefined,
    undefined,
    ts.createArrowFunction(
      undefined,
      undefined,
      undefined,
      undefined,
      undefined,
      ts.createBlock([ts.createStatement(promiseHandler)], true)
    )
  )
}

function createFetcher(meta: InputMeta) {
  return ts.createProperty(
    [
      ts.createDecorator(
        ts.createCall(ts.createIdentifier(Decorators.Intent), undefined, [ts.createLiteral(meta.intentName)])
      )
    ],
    undefined,
    meta.intentMethodName,
    undefined,
    undefined,
    undefined
  )
}

function createRefIdWatcher(meta: InputMeta) {
  const newValueName = 'newValueName'
  return ts.createMethod(
    [
      ts.createDecorator(
        ts.createCall(ts.createIdentifier(Decorators.Watch), undefined, [ts.createLiteral(meta.propName)])
      )
    ],
    undefined,
    undefined,
    meta.watcherName,
    undefined,
    undefined,
    [
      ts.createParameter(
        undefined,
        undefined,
        undefined,
        newValueName,
        undefined,
        ts.createKeywordTypeNode(ts.SyntaxKind.StringKeyword),
        undefined
      )
    ],
    undefined,
    ts.createBlock([
      ts.createIf(
        ts.createIdentifier(newValueName),
        ts.createBlock([
          ts.createStatement(
            ts.createCall(ts.createPropertyAccess(ts.createThis(), meta.loadMethodName), undefined, undefined)
          )
        ])
      )
    ])
  )
}

// // Nive to have
// function createLoadingProp(meta: InputMeta){

// }

type InputMeta = {
  propDeclarationName: string
  scope: string
  propName: string
  eventName: string
  intentName: string
  autoUpdate: boolean
  typeIdentifier?: ts.TypeNode
  intializer?: ts.Expression
  loadMethodName: string
  intentMethodName: string
  watcherName: string
}
