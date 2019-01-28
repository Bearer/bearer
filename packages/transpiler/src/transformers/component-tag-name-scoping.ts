import * as ts from 'typescript'

import Metadata from '../metadata'
import { Decorators } from '../constants'
import {
  decoratorNamed,
  getExpressionFromDecorator,
  hasDecoratorNamed,
  getDecoratorProperties
} from '../helpers/decorator-helpers'
import debug from '../logger'
const logger = debug.extend('tag-name-scoping')

const TAG = 'tag'

type TransformerOptions = {
  verbose?: true
  metadata: Metadata
}

export default function componentTagNameScoping({
  metadata
}: TransformerOptions): ts.TransformerFactory<ts.SourceFile> {
  function getNewTagName(tagName: string) {
    const component = metadata.components.find(aComponent => aComponent.initialTagName === tagName)
    if (!component) {
      return tagName
    }
    logger('rewriting %s to %s', tagName, component.finalTagName)
    return component.finalTagName
  }

  function visitJsxElement(node: ts.JsxElement) {
    const { openingElement, closingElement, children } = node

    const finalTagName = ts.createIdentifier(getNewTagName(stringFromTagName(openingElement)))
    return ts.updateJsxElement(
      node,
      ts.updateJsxOpeningElement(openingElement, finalTagName, openingElement.typeArguments, openingElement.attributes),
      children,
      ts.updateJsxClosingElement(closingElement, finalTagName)
    )
  }

  function visitJsxSelfClosingElement(node: ts.JsxSelfClosingElement) {
    return ts.updateJsxSelfClosingElement(
      node,
      ts.createIdentifier(getNewTagName(stringFromTagName(node))),
      node.typeArguments,
      node.attributes
    )
  }

  function visitDecoratedClassElement(node: ts.ClassDeclaration) {
    const decorators = node.decorators.map(decorator => {
      if (decoratorNamed(decorator, Decorators.Component)) {
        const finalTagName = getNewTagName((getExpressionFromDecorator(decorator, TAG) as any).text)
        const otherProperties = getDecoratorProperties(decorator, 0).properties.filter(
          (ol: ts.ObjectLiteralElement) => {
            return (ol.name as any).escapedText !== TAG
          }
        )
        return ts.updateDecorator(
          decorator,
          ts.createCall(ts.createIdentifier(Decorators.Component), undefined, [
            ts.createObjectLiteral(
              [ts.createPropertyAssignment(TAG, ts.createStringLiteral(finalTagName)), ...otherProperties],
              true
            )
          ])
        )
      }
      return decorator
    })
    return ts.updateClassDeclaration(
      node,
      decorators,
      node.modifiers,
      node.name,
      node.typeParameters,
      node.heritageClauses,
      node.members
    )
  }

  return _transformContext => {
    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      switch (tsNode.kind) {
        case ts.SyntaxKind.JsxElement:
          return ts.visitEachChild(visitJsxElement(tsNode as ts.JsxElement), visit, _transformContext)
        case ts.SyntaxKind.JsxSelfClosingElement:
          return ts.visitEachChild(
            visitJsxSelfClosingElement(tsNode as ts.JsxSelfClosingElement),
            visit,
            _transformContext
          )
        case ts.SyntaxKind.ClassDeclaration:
          if (hasDecoratorNamed(tsNode as ts.ClassDeclaration, Decorators.Component)) {
            const updatedNode = visitDecoratedClassElement(tsNode as ts.ClassDeclaration)
            return ts.visitEachChild(updatedNode, visit, _transformContext)
          }
        default:
      }
      return ts.visitEachChild(tsNode, visit, _transformContext)
    }

    return tsSourceFile => {
      if (tsSourceFile.isDeclarationFile) {
        return tsSourceFile
      }
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}

function stringFromTagName({ tagName }: { tagName: ts.JsxTagNameExpression }): string {
  return (tagName as ts.Identifier).escapedText.toString()
}
