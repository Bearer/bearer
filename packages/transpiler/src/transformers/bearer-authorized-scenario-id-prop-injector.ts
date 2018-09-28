/*
 * Append scenarioId prop to @Component containing bearer-authorized or navigator-auth-screen
 */
import * as ts from 'typescript'

import { Component, Env } from '../constants'
type TransformerOptions = {
  verbose?: true
}

const INJECTABLE = new Set([
  'bearer-authorized',
  'bearer-navigator-auth-screen'
])
export default function injectScenarioIdProp(
  _options: TransformerOptions = {}
): ts.TransformerFactory<ts.SourceFile> {
  function updateProperties(attributes: ts.JsxAttributes): ts.JsxAttributes {
    return ts.createJsxAttributes([
      ...attributes.properties.filter(
        p => (p.name as ts.Identifier).escapedText !== Component.scenarioId
      ),
      ts.createJsxAttribute(
        ts.createIdentifier(Component.scenarioId),
        ts.createLiteral(
          process.env[Env.BEARER_SCENARIO_ID] || Env.BEARER_SCENARIO_ID
        )
      )
    ])
  }

  function isInjectableJSXElement(
    element: ts.JsxSelfClosingElement | ts.JsxOpeningElement
  ): boolean {
    return INJECTABLE.has((element.tagName as ts.Identifier)
      .escapedText as string)
  }

  return _transformContext => {
    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      if (
        ts.isJsxSelfClosingElement(tsNode) &&
        isInjectableJSXElement(tsNode)
      ) {
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
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
