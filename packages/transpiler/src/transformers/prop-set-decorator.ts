/*
 * Prop Set Decorator
 * Applications like React do not expect props to dynamically update
 * This adds an event emmiter that emits '<finalTagName>-prop-set' events
 * with the data the attrubute was changed to. Allowing react etc to update
 * state
 *
 * NB: This must be included after event scoping and input/output so that it can work correctly
 */
import * as ts from 'typescript'

import { Decorators, Types } from '../constants'
import { extractBooleanOptions, getDecoratorNamed } from '../helpers/decorator-helpers'
import { getNodeName } from '../helpers/node-helpers'
import { TransformerOptions } from '../types'

import { ensureImportsFromCore } from './bearer'
import debug from '../logger'

const logger = debug.extend('prop-set-decorator')

const NEW_VALUE = 'newValue'
const PROP_SET_EMITTER_NAME = 'propSetEmitter'
const PROP_SET_EVENT_SUFFIX = '-prop-set'
const PROP_SET_WATCHER_SUFFIX = 'PropSetWatcher'

type TMutablePropMeta = {
  name: string
  tagName: string
}
export default function propSetDecorator(options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    const mutablePropMeta: TMutablePropMeta[] = []

    function visit(tagName: string) {
      return (tsNode: ts.Node): ts.VisitResult<ts.Node> => {
        if (ts.isPropertyDeclaration(tsNode)) {
          const decorator = getDecoratorNamed(tsNode, Decorators.Prop)
          if (decorator && isMutableProp(decorator)) {
            mutablePropMeta.push({ tagName, name: getNodeName(tsNode) })
          }
        }
        return ts.visitEachChild(tsNode, visit(tagName), _transformContext)
      }
    }

    const updateClass: ts.Visitor = (node: ts.Node) => {
      if (ts.isClassDeclaration(node)) {
        return ts.updateClassDeclaration(
          node,
          node.decorators,
          node.modifiers,
          node.name,
          node.typeParameters,
          node.heritageClauses,
          [...node.members, createPropSetEmitter(mutablePropMeta[0].tagName), ...createWatchers(mutablePropMeta)]
        )
      }

      return ts.visitEachChild(node, updateClass, _transformContext)
    }

    function isMutableProp(node: ts.Node) {
      return ts.isDecorator(node) && isMutableDecorator(node as ts.Decorator)
    }

    function isMutableDecorator(decorator: ts.Decorator) {
      const visitor = (node: ts.Node) => {
        if (ts.isObjectLiteralExpression(node)) {
          if (node && node.properties.length) {
            const { mutable } = extractBooleanOptions<{ mutable: boolean }>(node, ['mutable'])
            return mutable
          }
        }
      }
      return ts.forEachChild(decorator.expression, visitor)
    }
    return (sourceFile: ts.SourceFile): ts.SourceFile => {
      logger('processing %s', sourceFile.fileName)
      const meta = options.metadata.findComponentFrom(sourceFile)
      // only run this on root components
      if (meta && meta.isRoot) {
        ts.visitNode(sourceFile, visit(meta.finalTagName))
      }

      if (mutablePropMeta.length === 0) {
        // nothing to do
        return sourceFile
      }

      const sourceFileWithImports = ensureImportsFromCore(sourceFile, [
        Decorators.Watch,
        Decorators.Event,
        Types.EventEmitter
      ])
      return ts.visitNode(sourceFileWithImports, updateClass)
    }
  }
}

function createWatchers(props: TMutablePropMeta[]): ts.MethodDeclaration[] {
  return props.map(meta => {
    return ts.createMethod(
      [
        ts.createDecorator(
          ts.createCall(ts.createIdentifier(Decorators.Watch), undefined, [ts.createStringLiteral(meta.name)])
        )
      ],
      undefined,
      undefined,
      propSetWatcherName(meta),
      undefined,
      undefined,
      [ts.createParameter(undefined, undefined, undefined, NEW_VALUE, undefined, undefined, undefined)], // parameters
      undefined,
      ts.createBlock([ts.createStatement(createEventListener(meta))], true)
    )
  })
}
function createPropSetEmitter(tagName: string): ts.PropertyDeclaration {
  return ts.createProperty(
    [
      ts.createDecorator(
        ts.createCall(ts.createIdentifier(Decorators.Event), undefined, [
          ts.createObjectLiteral([
            ts.createPropertyAssignment('eventName', ts.createStringLiteral(propSetEventName(tagName))),
            ts.createPropertyAssignment('bubbles', ts.createTrue()),
            ts.createPropertyAssignment('cancelable', ts.createTrue())
          ])
        ])
      )
    ],
    undefined,
    PROP_SET_EMITTER_NAME,
    undefined,
    ts.createTypeReferenceNode(ts.createIdentifier(Types.EventEmitter), [
      ts.createTypeLiteralNode([
        ts.createPropertySignature(
          undefined,
          PROP_SET_EMITTER_NAME,
          undefined,
          ts.createKeywordTypeNode(ts.SyntaxKind.StringKeyword),
          undefined
        ),
        undefined
      ])
    ]),
    undefined
  )
}

function createEventListener(meta: TMutablePropMeta) {
  return ts.createCall(
    ts.createPropertyAccess(ts.createPropertyAccess(ts.createThis(), PROP_SET_EMITTER_NAME), 'emit'),
    undefined,
    [payload(meta)]
  )
}

function payload(meta: TMutablePropMeta) {
  return ts.createObjectLiteral([ts.createPropertyAssignment(meta.name, ts.createIdentifier(NEW_VALUE))])
}

const propSetWatcherName = (meta: TMutablePropMeta) => meta.name + PROP_SET_WATCHER_SUFFIX

const propSetEventName = (tagName: string) => tagName + PROP_SET_EVENT_SUFFIX
