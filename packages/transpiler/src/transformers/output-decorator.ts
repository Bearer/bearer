/*
 * Transformer boilerplate
 */
import * as ts from 'typescript'

import { Decorators, Types } from '../constants'
import { hasDecoratorNamed } from '../helpers/decorator-helpers'
import { getNodeName } from '../helpers/node-helpers'
import { TransformerOptions } from '../types'

import { ensureImportsFromCore } from './bearer'

const newValue = 'newValue'
const data = 'data'
const referenceId = 'referenceId'

export default function OutputDecorator(_options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    return tsSourceFile => {
      if (tsSourceFile.isDeclarationFile) {
        return tsSourceFile
      }

      const outputsMeta = retrieveOutputsMetas(tsSourceFile)

      if (!outputsMeta.length) {
        return tsSourceFile
      }

      const sourceFileWithImports = ensureImportsFromCore(tsSourceFile, [
        Types.BearerFetch,
        Types.EventEmitter,
        Decorators.Event,
        Decorators.Intent,
        Decorators.State,
        Decorators.Watch
      ])

      return ts.visitEachChild(sourceFileWithImports, visit(outputsMeta), _transformContext)
    }

    function retrieveOutputsMetas(tsSourceFile: ts.SourceFile): Array<OutputMeta> {
      const outputs: Array<OutputMeta> = []

      const visitor = (tsNode: ts.Node) => {
        if (ts.isPropertyDeclaration(tsNode) && hasDecoratorNamed(tsNode, Decorators.Output)) {
          const name = getNodeName(tsNode)
          outputs.push({
            emitMethodName: outputEventName(name), // TODO: retrieve from options
            intentName: `save${capitalize(name)}`, // TODO: retrieve from options
            intentPropertyName: name, // TODO: retrieve from options
            propDeclarationName: name,
            typeIdentifier: tsNode.type,
            initializer: tsNode.initializer,
            watchedPropName: name // TODO: retrieve from options
          })
        }
        return ts.visitEachChild(tsNode, visitor, _transformContext)
      }

      ts.visitEachChild(tsSourceFile, visitor, _transformContext)

      return outputs
    }

    function visit(outputsMeta: Array<OutputMeta>) {
      return (tsNode: ts.Node): ts.VisitResult<ts.Node> => {
        if (ts.isClassDeclaration(tsNode)) {
          return injectOuputStatements(tsNode, outputsMeta)
        }
        return ts.visitEachChild(tsNode, visit(outputsMeta), _transformContext)
      }
    }
  }
}
function injectOuputStatements(tsClass: ts.ClassDeclaration, outputsMeta: Array<OutputMeta>): ts.ClassDeclaration {
  const newMembers = outputsMeta.reduce(
    (members, meta) => {
      const inputMembers = [createIntent(meta), createEvent(meta), createState(meta), createWatcher(meta)]
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

function createEvent(meta: OutputMeta): ts.PropertyDeclaration {
  return ts.createProperty(
    [
      ts.createDecorator(
        ts.createCall(ts.createIdentifier(Decorators.Event), undefined, [ts.createStringLiteral(meta.watchedPropName)])
      )
    ],
    undefined,
    meta.emitMethodName,
    undefined,
    ts.createTypeReferenceNode(ts.createIdentifier(Types.EventEmitter), [
      ts.createKeywordTypeNode(ts.SyntaxKind.AnyKeyword)
    ]),
    undefined
  )
}

function createState(meta: OutputMeta): ts.PropertyDeclaration {
  return ts.createProperty(
    [ts.createDecorator(ts.createCall(ts.createIdentifier(Decorators.State), undefined, []))],
    undefined,
    meta.propDeclarationName,
    undefined,
    meta.typeIdentifier,
    meta.initializer
  )
}

function createWatcher(meta: OutputMeta): ts.MethodDeclaration {
  return ts.createMethod(
    [
      ts.createDecorator(
        ts.createCall(ts.createIdentifier(Decorators.Watch), undefined, [ts.createStringLiteral(meta.watchedPropName)])
      )
    ],
    undefined,
    undefined,
    `${meta.emitMethodName}Watcher`,
    undefined,
    undefined,
    [ts.createParameter(undefined, undefined, undefined, newValue, undefined, undefined, undefined)], // parameters
    undefined,
    ts.createBlock(
      [
        ts.createIf(
          ts.createIdentifier(newValue),
          ts.createBlock([ts.createStatement(createIntentCall(meta))], true),
          ts.createBlock([
            ts.createStatement(
              createEmitCall(meta, [
                ts.createPropertyAssignment(meta.propDeclarationName, ts.createIdentifier(newValue))
              ])
            )
          ])
        )
      ],
      true
    )
  )
}

function createIntent(meta: OutputMeta): ts.PropertyDeclaration {
  return ts.createProperty(
    [
      ts.createDecorator(
        ts.createCall(ts.createIdentifier(Decorators.Intent), undefined, [ts.createStringLiteral(meta.intentName)])
      )
    ],
    undefined,
    meta.intentName,
    undefined,
    ts.createTypeReferenceNode(ts.createIdentifier(Types.BearerFetch), [
      meta.typeIdentifier || ts.createKeywordTypeNode(ts.SyntaxKind.AnyKeyword)
    ]),
    undefined
  )
}

function createIntentCall(meta: OutputMeta) {
  const newValueInitializer = ts.createBinary(
    ts.createIdentifier(data),
    ts.createToken(ts.SyntaxKind.BarBarToken),
    ts.createPropertyAccess(ts.createThis(), meta.propDeclarationName)
  )

  const emit = createEmitCall(meta, [
    ts.createShorthandPropertyAssignment(referenceId),
    ts.createPropertyAssignment(meta.propDeclarationName, newValueInitializer)
  ])

  return ts.createCall(
    ts.createPropertyAccess(
      ts.createCall(ts.createPropertyAccess(ts.createThis(), meta.intentName), undefined, [
        ts.createObjectLiteral([
          ts.createPropertyAssignment(
            'body',
            ts.createObjectLiteral([
              ts.createPropertyAssignment(meta.intentPropertyName, ts.createIdentifier(newValue))
            ])
          )
        ])
      ]),
      'then'
    ),
    undefined,
    [
      ts.createArrowFunction(
        undefined,
        undefined,
        // [],
        [
          ts.createParameter(
            undefined,
            undefined,
            undefined,
            ts.createObjectBindingPattern([
              ts.createBindingElement(undefined, undefined, referenceId, undefined),
              ts.createBindingElement(undefined, undefined, data, undefined)
            ]),
            undefined,
            undefined,
            undefined
          )
        ],
        undefined,
        ts.createToken(ts.SyntaxKind.EqualsGreaterThanToken),
        ts.createBlock([ts.createStatement(emit)], true)
      )
    ]
  )
}

function createEmitCall(meta: OutputMeta, properties: Array<ts.ObjectLiteralElementLike>): ts.CallExpression {
  return ts.createCall(
    ts.createPropertyAccess(ts.createPropertyAccess(ts.createThis(), meta.emitMethodName), 'emit'),
    undefined,
    [ts.createObjectLiteral(properties)]
  )
}

export function outputEventName(prefix: string, suffix?: string): string {
  const _suffix = suffix || 'Saved'
  return `${prefix}${capitalize(_suffix)}`
}

function capitalize(string: string): string {
  return string.charAt(0).toUpperCase() + string.slice(1)
}

type OutputMeta = {
  emitMethodName: string
  intentName: string
  propDeclarationName: string
  initializer: ts.Expression
  typeIdentifier?: ts.TypeNode
  watchedPropName: string
  intentPropertyName: string
}
