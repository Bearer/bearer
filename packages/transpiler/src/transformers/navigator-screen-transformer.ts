/*
 * Rewrite navigator-screen if they do not use renderFunc
 */
import * as ts from 'typescript'

import { TransformerOptions } from '../types'

const NAVIGATOR_SCREEN_TAG_NAME = 'bearer-navigator-screen'

export default function PropImporter({  }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  function moveSlotToRenderFuncProp(jsxNode: ts.JsxElement): ts.JsxSelfClosingElement {
    return ts.createJsxSelfClosingElement(
      jsxNode.openingElement.tagName,
      jsxNode.openingElement.typeArguments,
      ts.createJsxAttributes([
        ...jsxNode.openingElement.attributes.properties,
        ts.createJsxAttribute(
          ts.createIdentifier('renderFunc'),
          ts.createJsxExpression(
            undefined,
            ts.createArrowFunction(
              undefined /* modifiers */,
              undefined,
              undefined,
              undefined,
              ts.createToken(ts.SyntaxKind.EqualsGreaterThanToken),
              ts.createArrayLiteral([
                ...jsxNode.children.filter(child => {
                  // string not scoped are not supported because they are transformed into => ()
                  return ts.isJsxElement(child) || ts.isJsxSelfClosingElement(child)
                })
              ] as Array<ts.Expression>)
            )
          )
        )
      ])
    )
  }

  return _transformContext => {
    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      if (ts.isJsxElement(tsNode)) {
        if ((tsNode.openingElement.tagName as ts.Identifier).escapedText === NAVIGATOR_SCREEN_TAG_NAME) {
          return ts.visitEachChild(moveSlotToRenderFuncProp(tsNode), visit, _transformContext)
        }
      }
      return ts.visitEachChild(tsNode, visit, _transformContext)
    }

    return tsSourceFile => {
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
