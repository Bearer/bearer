/*
 * Transformer boilerplate
 */
import { TOutputDecoratorOptions } from '@bearer/types/lib/input-output-decorators'
import * as ts from 'typescript'

import { Decorators, Properties, Types } from '../constants'
import {
  extractBooleanOptions,
  extractStringOptions,
  getDecoratorNamed,
  extractArrayOptions
} from '../helpers/decorator-helpers'
import { addAutoLoad, createFetcher, createLoadResourceMethod } from '../helpers/generator-helpers'
import { initialName, retrieveFetcherName, retrieveIntentName, loadName } from '../helpers/name-helpers'
import { getNodeName } from '../helpers/node-helpers'
import { capitalize } from '../helpers/string'
import { TCreateLoadResourceMethod, TransformerOptions, OutputMeta, InputMeta } from '../types'

import { ensureImportsFromCore } from './bearer'
import { retrieveInputsMetas } from './input-decorator'

const newValue = 'newValue'
const data = 'data'

export default function outputDecorator({ metadata }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    return tsSourceFile => {
      if (tsSourceFile.isDeclarationFile) {
        return tsSourceFile
      }

      const outputsMeta = retrieveOutputsMetas(tsSourceFile, _transformContext)

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

      return ts.visitEachChild(sourceFileWithImports, visit(outputsMeta, tsSourceFile), _transformContext)
    }

    function visit(outputsMeta: OutputMeta[], tsSourceFile: ts.SourceFile) {
      return (tsNode: ts.Node): ts.VisitResult<ts.Node> => {
        if (ts.isClassDeclaration(tsNode)) {
          return injectOuputStatements(
            tsNode,
            outputsMeta,
            retrieveInputsMetas(tsSourceFile, metadata, _transformContext)
          )
        }
        return ts.visitEachChild(tsNode, visit(outputsMeta, tsSourceFile), _transformContext)
      }
    }
  }
}

function injectOuputStatements(
  tsClass: ts.ClassDeclaration,
  outputsMeta: OutputMeta[],
  inputsMeta: InputMeta[]
): ts.ClassDeclaration {
  const classNode = outputsMeta.reduce((classDeclaration, meta) => addAutoLoad(classDeclaration, meta), tsClass)
  const newMembers = outputsMeta.reduce(
    (members, meta) => {
      const outputMembers = [
        createIntent(meta),
        createEvent(meta),
        ...createStates(meta),
        createProp(meta),
        ...createWatchers(meta),
        createInitialFetcher(meta),
        createLoadResourceMethod(
          {
            ...(meta as TCreateLoadResourceMethod),
            propDeclarationName: initialName(meta.propDeclarationName)
          },
          [...outputsMeta, ...inputsMeta]
        )
      ]
      return members.concat(outputMembers)
    },
    [...classNode.members]
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

// This would let developer use passed reference
// @Prop({ mutable: true }) propDeclarationNameRefId: string
function createProp(meta: OutputMeta): ts.PropertyDeclaration {
  return ts.createProperty(
    [
      ts.createDecorator(
        ts.createCall(ts.createIdentifier(Decorators.Prop), undefined, [
          ts.createObjectLiteral([ts.createPropertyAssignment(ts.createLiteral('mutable'), ts.createTrue())])
        ])
      )
    ],
    undefined,
    meta.propDeclarationNameRefId,
    undefined,
    ts.createKeywordTypeNode(ts.SyntaxKind.StringKeyword),
    undefined
  )
}

function createStates(meta: OutputMeta): ts.PropertyDeclaration[] {
  return [
    // @State() propDeclarationNameInitial: Type = initiailizer
    createGenericState(meta, initialName),
    // @State() propDeclarationName: Type = initiailizer
    createGenericState(meta)
  ]
}

function createGenericState(meta: OutputMeta, nameTransformation = (x: string) => x): ts.PropertyDeclaration {
  // @State() propDeclarationName[suffix]: Type = initiailizer
  return ts.createProperty(
    [ts.createDecorator(ts.createCall(ts.createIdentifier(Decorators.State), undefined, []))],
    undefined,
    nameTransformation(meta.propDeclarationName),
    undefined,
    meta.typeIdentifier,
    meta.initializer
  )
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
      ts.createBlock(
        [
          ts.createStatement(createIntentCall(meta)),
          ts.createStatement(
            ts.createBinary(
              ts.createPropertyAccess(ts.createThis(), initialName(meta.propDeclarationName)),
              ts.SyntaxKind.EqualsToken,
              ts.createNull()
            )
          )
        ],
        true
      )
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
            ts.createPropertyAccess(ts.createThis(), meta.intentReferenceIdValue || meta.propDeclarationNameRefId)
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

function createEmitCall(meta: OutputMeta, properties: ts.ObjectLiteralElementLike[]): ts.CallExpression {
  return ts.createCall(
    ts.createPropertyAccess(ts.createPropertyAccess(ts.createThis(), meta.eventName), 'emit'),
    undefined,
    [ts.createObjectLiteral(properties)]
  )
}

export function refIdName(name: string): string {
  return `${name}Id`
}

export function saveIntentName(name: string): string {
  return `save${capitalize(name)}`
}

export function outputEventName(prefix: string, suffix?: string): string {
  const _suffix = suffix || 'Saved'
  return `${prefix}${capitalize(_suffix)}`
}

function createInitialFetcher(meta) {
  const metaForInitial = {
    intentName: retrieveIntentName(meta.propDeclarationName)
  }
  return createFetcher({ ...meta, ...metaForInitial })
}

export function retrieveOutputsMetas(
  tsSourceFile: ts.SourceFile,
  _transformContext: ts.TransformationContext
): OutputMeta[] {
  const outputs: OutputMeta[] = []

  const visitor = (tsNode: ts.Node) => {
    if (ts.isPropertyDeclaration(tsNode)) {
      const decorator = getDecoratorNamed(tsNode, Decorators.Output)
      if (decorator) {
        const name = getNodeName(tsNode)
        const callArgs = (decorator.expression as ts.CallExpression).arguments[0] as ts.ObjectLiteralExpression
        const options = !callArgs
          ? {}
          : {
              ...extractStringOptions<TOutputDecoratorOptions>(callArgs, [
                'eventName',
                'intentName',
                'intentPropertyName',
                'propertyWatchedName',
                'referenceKeyName',
                'intentReferenceIdValue',
                'intentReferenceIdKeyName'
              ]),
              ...extractArrayOptions<{ intentArguments: string[] }>(callArgs, ['intentArguments']),
              ...extractBooleanOptions<TOutputDecoratorOptions>(callArgs, ['autoLoad'])
            }
        outputs.push({
          eventName: outputEventName(name),
          intentName: saveIntentName(name),
          intentMethodName: retrieveFetcherName(name),
          intentPropertyName: name,
          propDeclarationName: name,
          propDeclarationNameRefId: refIdName(name),
          loadMethodName: loadName(name),
          intentReferenceIdKeyName: Properties.ReferenceId,
          typeIdentifier: tsNode.type,
          initializer: tsNode.initializer,
          referenceKeyName: Properties.ReferenceId,
          propertyWatchedName: name,
          propertyReferenceIdName: refIdName(name),
          autoLoad: true,
          intentArguments: [],
          ...options
        })
      }
    }
    return ts.visitEachChild(tsNode, visitor, _transformContext)
  }

  ts.visitEachChild(tsSourceFile, visitor, _transformContext)

  return outputs
}
