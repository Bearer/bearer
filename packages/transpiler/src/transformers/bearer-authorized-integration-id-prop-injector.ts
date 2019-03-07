/*
 * Append integrationId prop to @Component containing bearer-authorized or navigator-auth-screen
 */
import * as ts from 'typescript'

import { Component, Env } from '../constants'
import debug from '../logger'

const logger = debug.extend('inject-integration-id-prop')

type TransformerOptions = {
  verbose?: true
}

const INJECTABLE = new Set(['bearer-authorized', 'bearer-navigator-auth-screen'])

export default function injectIntegrationIdProp(
  _options: TransformerOptions = {}
): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      if (ts.isJsxSelfClosingElement(tsNode) && isInjectableJSXElement(tsNode)) {
        return ts.updateJsxSelfClosingElement(
          tsNode,
          tsNode.tagName,
          tsNode.typeArguments,
          updateProperties(tsNode.attributes)
        )
      }
      if (ts.isJsxOpeningElement(tsNode) && isInjectableJSXElement(tsNode)) {
        const updatedNode = ts.updateJsxOpeningElement(
          tsNode,
          tsNode.tagName,
          tsNode.typeArguments,
          updateProperties(tsNode.attributes)
        )
        return ts.visitEachChild(updatedNode, visit, _transformContext)
      }
      return ts.visitEachChild(tsNode, visit, _transformContext)
    }

    return tsSourceFile => {
      logger('processing %s', tsSourceFile.fileName)
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}

function updateProperties(attributes: ts.JsxAttributes): ts.JsxAttributes {
  return ts.createJsxAttributes([
    ...attributes.properties.filter(p => (p.name as ts.Identifier).escapedText !== Component.integrationId),
    ts.createJsxAttribute(
      ts.createIdentifier(Component.integrationId),
      ts.createLiteral(process.env[Env.BEARER_INTEGRATION_ID] || Env.BEARER_INTEGRATION_ID)
    )
  ])
}

function isInjectableJSXElement(element: ts.JsxSelfClosingElement | ts.JsxOpeningElement): boolean {
  return INJECTABLE.has((element.tagName as ts.Identifier).escapedText as string)
}
