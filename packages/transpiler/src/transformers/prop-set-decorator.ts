/*
 * Transformer boilerplate
 */
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { getDecoratorNamed } from '../helpers/decorator-helpers'
import { getNodeName } from '../helpers/node-helpers'
import { TransformerOptions } from '../types'

import { ensureImportsFromCore } from './bearer'

const NEW_VALUE = 'newValue'
type TMutablePropMeta = {
  name: string
  tagName: string
}
export default function PropSetDecorator(options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    const mutablePropMeta: Array<TMutablePropMeta> = []

    function visit(tagName: string) {
      return (tsNode: ts.Node): ts.VisitResult<ts.Node> => {
        if (ts.isPropertyDeclaration(tsNode)) {
          const decorator = getDecoratorNamed(tsNode, Decorators.Prop)
          if (decorator && isMutableProp(decorator)) {
            mutablePropMeta.push({ name: getNodeName(tsNode), tagName })
          }
        }
        return ts.visitEachChild(tsNode, visit(tagName), _transformContext)
      }
    }

    const updateClass: ts.Visitor = (node: ts.Node) => {
      if (ts.isClassDeclaration(node)) {
        const prop = node as ts.ClassDeclaration
        return ts.updateClassDeclaration(
          prop,
          prop.decorators,
          prop.modifiers,
          prop.name,
          prop.typeParameters,
          prop.heritageClauses,
          [...prop.members, ...createWatchers(mutablePropMeta)]
        )
      }

      return ts.visitEachChild(node, updateClass, _transformContext)
    }

    return (sourceFile: ts.SourceFile): ts.SourceFile => {
      const meta = options.metadata.findComponentFrom(sourceFile)
      //only run this on root components
      if (meta && meta.isRoot) {
        ts.visitNode(sourceFile, visit(meta.finalTagName))
      }

      if (mutablePropMeta.length) {
        const sourceFileWithImports = ensureImportsFromCore(sourceFile, [Decorators.Watch])
        return ts.visitNode(sourceFileWithImports, updateClass)
      }
      return sourceFile
    }
  }
}

function hasBooleanProperty(
  node: ts.ObjectLiteralExpression,
  prop: string,
  kind: ts.SyntaxKind.TrueKeyword | ts.SyntaxKind.FalseKeyword
) {
  const argument = node.properties[0]
  let identifier: any = { escapedText: 'undefined' }
  let initializer: any = { kind: ts.SyntaxKind.UndefinedKeyword }
  if (argument.kind === ts.SyntaxKind.PropertyAssignment) {
    identifier = argument.name as ts.Identifier
    initializer = argument.initializer
  }
  return (identifier.escapedText || identifier.text) === prop && initializer.kind === kind
}

function isMutableDecorator(decorator: ts.Decorator) {
  let isMutable = false
  const visitor = (node: ts.Node): ts.VisitResult<ts.Node> => {
    if (ts.isObjectLiteralExpression(node)) {
      isMutable = hasBooleanProperty(node, 'mutable', ts.SyntaxKind.TrueKeyword)
    }

    return node
  }
  ts.visitNodes((decorator.expression as ts.CallExpression).arguments, visitor)

  return isMutable
}

function isMutableProp(node: ts.Node) {
  return ts.isDecorator(node) && isMutableDecorator(node as ts.Decorator)
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
      `${meta.name}PropSetWatcher`,
      undefined,
      undefined,
      [ts.createParameter(undefined, undefined, undefined, NEW_VALUE, undefined, undefined, undefined)], // parameters
      undefined,
      ts.createBlock([ts.createStatement(createEventListener(meta))], true)
    )
  })
}

function createEventListener(meta: TMutablePropMeta) {
  const preparedEvent = ts.createNew(ts.createIdentifier('CustomEvent'), undefined, [
    ts.createStringLiteral(`${meta.tagName}-prop-set`),
    payload(meta)
  ])

  return ts.createCall(ts.createPropertyAccess(ts.createIdentifier('document'), 'dispatchEvent'), undefined, [
    preparedEvent
  ])
}

function payload(meta: TMutablePropMeta) {
  const detail = ts.createObjectLiteral([ts.createPropertyAssignment(meta.name, ts.createIdentifier(NEW_VALUE))])
  return ts.createObjectLiteral([ts.createPropertyAssignment('detail', detail)])
}
