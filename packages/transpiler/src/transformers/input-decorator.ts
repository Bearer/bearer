/*
 * Input Transformer
 */
import { TInputDecoratorOptions } from '@bearer/types/lib/input-output-decorators'
import * as ts from 'typescript'

import { Decorators, Properties } from '../constants'
import {
  extractBooleanOptions,
  extractStringOptions,
  getDecoratorNamed,
  extractArrayOptions
} from '../helpers/decorator-helpers'
import { addAutoLoad, createFetcher, createLoadDataCall, createLoadResourceMethod } from '../helpers/generator-helpers'
import { loadName as _loadName, retrieveFetcherName, retrieveIntentName } from '../helpers/name-helpers'
import { getNodeName } from '../helpers/node-helpers'
import { capitalize } from '../helpers/string'
import { InputMeta, TransformerOptions, OutputMeta } from '../types'

import { ensureImportsFromCore } from './bearer'
import { outputEventName, refIdName, retrieveOutputsMetas } from './output-decorator'
import Metadata from '../metadata'

export default function inputDecorator({ metadata }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    return tsSourceFile => {
      if (tsSourceFile.isDeclarationFile) {
        return tsSourceFile
      }
      // gather inputs information
      const inputsMeta = retrieveInputsMetas(tsSourceFile, metadata, _transformContext)

      if (!inputsMeta.length) {
        return tsSourceFile
      }

      const sourceFileWithImports = ensureImportsFromCore(tsSourceFile, [
        Decorators.Listen,
        Decorators.State,
        Decorators.Intent,
        Decorators.Watch,
        Decorators.Prop
      ])

      return ts.visitEachChild(sourceFileWithImports, replaceInputVisitor(inputsMeta, tsSourceFile), _transformContext)
    }

    function injectInputStatements(
      tsClass: ts.ClassDeclaration,
      inputsMeta: InputMeta[],
      outputsMeta: OutputMeta[]
    ): ts.ClassDeclaration {
      const classNode = inputsMeta.reduce((classDeclaration, meta) => addAutoLoad(classDeclaration, meta), tsClass)
      const newMembers = inputsMeta.reduce(
        (members, meta) => {
          // create @State()
          const inputMembers = [
            createLocalStateProperty(meta),
            createEventListener(meta),
            createLoadResourceMethod(meta, [...inputsMeta, ...outputsMeta]),
            createFetcher(meta),
            createRefIdWatcher(meta)
          ]
          if (
            !classNode.members.find(m => ts.isPropertyDeclaration(m) && getNodeName(m) === meta.propertyReferenceIdName)
          ) {
            inputMembers.push(createRefIdProp(meta))
          }
          return members.concat(inputMembers)
        },
        [...classNode.members]
      )

      return ts.updateClassDeclaration(
        classNode,
        classNode.decorators,
        classNode.modifiers,
        classNode.name,
        classNode.typeParameters,
        classNode.heritageClauses,
        newMembers
      )
    }

    function replaceInputVisitor(
      inputsMeta: InputMeta[],
      tsSourceFile: ts.SourceFile
    ): (tsNode: ts.Node) => ts.VisitResult<ts.Node> {
      return (tsNode: ts.Node) => {
        if (ts.isClassDeclaration(tsNode)) {
          return ts.visitEachChild(
            injectInputStatements(tsNode, inputsMeta, retrieveOutputsMetas(tsSourceFile, _transformContext)),
            replaceInputVisitor(inputsMeta, tsSourceFile),
            _transformContext
          )
        }
        return ts.visitEachChild(tsNode, replaceInputVisitor(inputsMeta, tsSourceFile), _transformContext)
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
    ts.createIdentifier(meta.propertyReferenceIdName),
    undefined,
    ts.createKeywordTypeNode(ts.SyntaxKind.StringKeyword),
    undefined
  )
}

// This will make the prop change when an event is triggered
function createEventListener(meta: InputMeta) {
  const referenceIdIdentifier = ts.createIdentifier('event.detail.referenceId')
  const propAccessIdentifier = ts.createPropertyAccess(ts.createThis(), meta.propertyReferenceIdName)
  return ts.createMethod(
    [
      ts.createDecorator(
        ts.createCall(ts.createIdentifier(Decorators.Listen), undefined, [
          ts.createLiteral(`body:${meta.group}:${meta.eventName}`)
        ])
      )
    ],
    undefined,
    undefined,
    `${meta.propertyReferenceIdName}Changed`,
    undefined,
    undefined,
    [ts.createParameter(undefined, undefined, undefined, ts.createIdentifier('event'), undefined, undefined)],
    undefined,
    ts.createBlock(
      [
        // if(event.detail.referenceId !== referenceIdIdentifier) {
        // this.referenceIdIdentifier = event.detail.referenceId
        // } else { this.loadData()}
        ts.createIf(
          ts.createBinary(propAccessIdentifier, ts.SyntaxKind.ExclamationEqualsEqualsToken, referenceIdIdentifier),
          ts.createBlock([
            // we don't have any scoping for now, this is a temporary workaround until we handle that in a better way
            // nothing gets updated if a reference has been provided
            ts.createIf(
              ts.createPrefix(ts.SyntaxKind.ExclamationToken, propAccessIdentifier),
              ts.createBlock([
                ts.createStatement(
                  ts.createBinary(propAccessIdentifier, ts.SyntaxKind.EqualsToken, referenceIdIdentifier)
                )
              ])
            )
          ]),
          ts.createBlock([createLoadDataCall(meta)])
        )
      ],
      true
    )
  )
}

function createRefIdWatcher(meta: InputMeta) {
  const newValueName = 'newValueName'
  return ts.createMethod(
    [
      ts.createDecorator(
        ts.createCall(ts.createIdentifier(Decorators.Watch), undefined, [
          ts.createLiteral(meta.propertyReferenceIdName)
        ])
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
    ts.createBlock(
      [ts.createIf(ts.createIdentifier(newValueName), ts.createBlock([createLoadDataCall(meta)], true))],
      true
    )
  )
}

// tslint:disable-next-line function-name
function _watchName(name: string): string {
  return `_watch${capitalize(name)}`
}

export function retrieveInputsMetas(
  sourcefile: ts.SourceFile,
  metadata: Metadata,
  _transformContext: ts.TransformationContext
): InputMeta[] {
  const inputs: InputMeta[] = []

  const visitor = (tsNode: ts.Node) => {
    if (ts.isPropertyDeclaration(tsNode)) {
      const decorator = getDecoratorNamed(tsNode, Decorators.Input)
      if (decorator) {
        const options = extractInputOptions(decorator)
        const name = getNodeName(tsNode)
        const component = metadata.findComponentFrom(sourcefile)
        inputs.push({
          propDeclarationName: name,
          group: component.group,
          propertyReferenceIdName: refIdName(name),
          eventName: outputEventName(name),
          intentName: retrieveIntentName(name),
          intentMethodName: retrieveFetcherName(name), // TODO: retrieve from options
          autoLoad: true,
          loadMethodName: _loadName(name),
          typeIdentifier: tsNode.type,
          intializer: tsNode.initializer,
          watcherName: _watchName(name),
          intentReferenceIdKeyName: Properties.ReferenceId,
          intentArguments: [],
          ...options
        })
      }
    }
    return ts.visitEachChild(tsNode, visitor, _transformContext)
  }

  ts.visitEachChild(sourcefile, visitor, _transformContext)

  return inputs
}

function extractInputOptions(decorator: ts.Decorator): Partial<TInputDecoratorOptions> {
  const callArgs = (decorator.expression as ts.CallExpression).arguments[0] as ts.ObjectLiteralExpression
  return !callArgs
    ? {}
    : {
        ...extractStringOptions<TInputDecoratorOptions>(callArgs, [
          'group',
          'eventName',
          'intentName',
          'propertyReferenceIdName',
          'intentReferenceIdKeyName'
        ]),
        ...extractArrayOptions<{ intentArguments: string[] }>(callArgs, ['intentArguments']),
        ...extractBooleanOptions<TInputDecoratorOptions>(callArgs, ['autoLoad'])
      }
}
