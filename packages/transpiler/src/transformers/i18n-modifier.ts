/*
 * scope i18n component and helpers
 */
import * as ts from 'typescript'

// At the same time it ensures we will have the accessor present ;-)
import { shouldProcessFile as hasAccessor } from './scenario-id-accessor-injector'

import { TransformerOptions } from '../types'
import { Component } from '../constants'

const I18N_TAG = 'bearer-i18n'

function visitJsxSelfClosingElement(node: ts.JsxSelfClosingElement): ts.JsxSelfClosingElement {
  if ((node.tagName as ts.Identifier).escapedText.toString().toLowerCase() !== I18N_TAG) {
    return node
  }
  return ts.updateJsxSelfClosingElement(
    node,
    node.tagName,
    node.typeArguments,
    ts.createJsxAttributes([
      ...node.attributes.properties,
      ts.createJsxAttribute(
        ts.createIdentifier('scope'),
        ts.createJsxExpression(undefined, ts.createPropertyAccess(ts.createThis(), Component.scenarioIdAccessor))
      )
    ])
  )
}

export default function i18nModifier(_options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      // if bearer-i18n tag add scope
      switch (tsNode.kind) {
        case ts.SyntaxKind.JsxSelfClosingElement:
          return ts.visitEachChild(
            visitJsxSelfClosingElement(tsNode as ts.JsxSelfClosingElement),
            visit,
            _transformContext
          )
      }
      return ts.visitEachChild(tsNode, visit, _transformContext)
    }

    return tsSourceFile => {
      if (!hasAccessor(tsSourceFile)) {
        return tsSourceFile
      }
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
