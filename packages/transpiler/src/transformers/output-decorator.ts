/*
 * Transformer boilerplate
 */
import { TOutputDecoratorOptions } from '@bearer/types/lib/input-output-decorators'
import * as ts from 'typescript'

import { Decorators, Properties, Types } from '../constants'
import { extractStringOptions, getDecoratorNamed } from '../helpers/decorator-helpers'
import { getNodeName } from '../helpers/node-helpers'
import { capitalize } from '../helpers/string'
import { TransformerOptions } from '../types'

import { ensureImportsFromCore } from './bearer'

const newValue = 'newValue'
const data = 'data'

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
        if (ts.isPropertyDeclaration(tsNode)) {
          const decorator = getDecoratorNamed(tsNode, Decorators.Output)
          if (decorator) {
            const name = getNodeName(tsNode)
            const callArgs = (decorator.expression as ts.CallExpression).arguments[0] as ts.ObjectLiteralExpression
            const options = !callArgs
              ? {}
              : extractStringOptions<TOutputDecoratorOptions>(callArgs, [
                  'eventName',
                  'intentName',
                  'intentPropertyName',
                  'propertyWatchedName',
                  'referenceKeyName'
                ])
            outputs.push({
              eventName: outputEventName(name),
              intentName: saveIntentName(name),
              intentPropertyName: name,
              propDeclarationName: name,
              propDeclarationNameRefId: refIdName(name),
              intentReferenceIdKeyName: Properties.ReferenceId,
              typeIdentifier: tsNode.type,
              initializer: tsNode.initializer,
              referenceKeyName: Properties.ReferenceId,
              propertyWatchedName: name,
              ...options
            })
          }
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
      const inputMembers = [createIntent(meta), createEvent(meta), ...createStates(meta), ...createWatchers(meta)]
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
  // @Event() propertyWatchedName: EventEmitter<any>;
  return ts.createProperty(
    [
      ts.createDecorator(
        ts.createCall(ts.createIdentifier(Decorators.Event), undefined, [
          ts.createStringLiteral(meta.propertyWatchedName)
        ])
      )
    ],
    undefined,
    meta.eventName,
    undefined,
    ts.createTypeReferenceNode(ts.createIdentifier(Types.EventEmitter), [
      ts.createTypeLiteralNode([
        ts.createPropertySignature(
          undefined,
          meta.referenceKeyName,
          undefined,
          ts.createKeywordTypeNode(ts.SyntaxKind.StringKeyword),
          undefined
        ),
        ts.createPropertySignature(undefined, meta.propDeclarationName, undefined, meta.typeIdentifier, undefined)
      ])
    ]),
    undefined
  )
}

function createStates(meta: OutputMeta): ts.PropertyDeclaration[] {
  return [
    // @State() propDeclarationName: Type = initiailizer
    ts.createProperty(
      [ts.createDecorator(ts.createCall(ts.createIdentifier(Decorators.State), undefined, []))],
      undefined,
      meta.propDeclarationName,
      undefined,
      meta.typeIdentifier,
      meta.initializer
    ),
    // @State() propDeclarationNameRefId: Type = initiailizer
    ts.createProperty(
      [ts.createDecorator(ts.createCall(ts.createIdentifier(Decorators.State), undefined, []))],
      undefined,
      meta.propDeclarationNameRefId,
      undefined,
      ts.createKeywordTypeNode(ts.SyntaxKind.StringKeyword),
      undefined
    )
  ]
}

function createWatchers(meta: OutputMeta): ts.MethodDeclaration[] {
  return [
    ts.createMethod(
      [
        ts.createDecorator(
          ts.createCall(ts.createIdentifier(Decorators.Watch), undefined, [
            ts.createStringLiteral(meta.propertyWatchedName)
          ])
        )
      ],
      undefined,
      undefined,
      `${meta.eventName}Watcher`,
      undefined,
      undefined,
      [ts.createParameter(undefined, undefined, undefined, newValue, undefined, undefined, undefined)], // parameters
      undefined,
      ts.createBlock([ts.createStatement(createIntentCall(meta))], true)
    )
  ]
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
  // TODO use referenceKeyName
  const emit = createEmitCall(meta, [
    ts.createShorthandPropertyAssignment(meta.referenceKeyName),
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
          ),
          ts.createPropertyAssignment(
            meta.intentReferenceIdKeyName,
            ts.createPropertyAccess(ts.createThis(), meta.propDeclarationNameRefId)
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
              ts.createBindingElement(
                undefined,
                meta.referenceKeyName !== Properties.ReferenceId ? Properties.ReferenceId : undefined,
                meta.referenceKeyName,
                undefined
              ),
              ts.createBindingElement(undefined, undefined, data, undefined)
            ]),
            undefined,
            undefined,
            undefined
          )
        ],
        undefined,
        ts.createToken(ts.SyntaxKind.EqualsGreaterThanToken),
        ts.createBlock(
          [
            ts.createStatement(emit),
            ts.createStatement(
              ts.createBinary(
                ts.createPropertyAccess(ts.createThis(), meta.propDeclarationNameRefId),
                ts.SyntaxKind.EqualsToken,
                ts.createIdentifier(meta.referenceKeyName)
              )
            )
          ],
          true
        )
      )
    ]
  )
}

function createEmitCall(meta: OutputMeta, properties: Array<ts.ObjectLiteralElementLike>): ts.CallExpression {
  return ts.createCall(
    ts.createPropertyAccess(ts.createPropertyAccess(ts.createThis(), meta.eventName), 'emit'),
    undefined,
    [ts.createObjectLiteral(properties)]
  )
}

export function refIdName(name: string): string {
  return `${name}RefId`
}

export function saveIntentName(name: string): string {
  return `save${capitalize(name)}`
}

export function outputEventName(prefix: string, suffix?: string): string {
  const _suffix = suffix || 'Saved'
  return `${prefix}${capitalize(_suffix)}`
}

type OutputMeta = TOutputDecoratorOptions & {
  propDeclarationName: string
  propDeclarationNameRefId: string
  initializer: ts.Expression
  typeIdentifier?: ts.TypeNode
}
