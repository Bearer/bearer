import * as ts from 'typescript'

import { Decorators, Properties } from '../constants'
import { decoratorNamed } from '../helpers/decorator-helpers'
import { isBearerEvent } from './event-name-scoping'
import * as Case from 'case'

export default function eventNameNormalizer(): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    function normalizeEventName(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      if (
        ts.isPropertyAssignment(tsNode) &&
        (tsNode.name as ts.Identifier).escapedText.toString() === Properties.eventName
      ) {
        const normalizedName = normalize((tsNode.initializer as ts.StringLiteral).text.toString())
        return ts.createPropertyAssignment(tsNode.name, ts.createLiteral(normalizedName))
      }
      return ts.visitEachChild(tsNode, normalizeEventName, _transformContext)
    }

    function normalizeListenerEventName(tsNode: ts.Decorator): ts.VisitResult<ts.Node> {
      const expression = tsNode.expression
      if (ts.isCallExpression(expression)) {
        const [eventName, ...rest] = expression.arguments
        const normalizedName = normalize((eventName as ts.StringLiteral).text)
        return ts.updateDecorator(
          tsNode,
          ts.createCall(expression.expression, expression.typeArguments, [ts.createLiteral(normalizedName), ...rest])
        )
      }
      return ts.visitEachChild(tsNode, normalizeEventName, _transformContext)
    }

    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      if (ts.isDecorator(tsNode) && decoratorNamed(tsNode, Decorators.Event)) {
        return normalizeEventName(tsNode)
      }

      if (ts.isDecorator(tsNode) && decoratorNamed(tsNode, Decorators.Listen)) {
        return normalizeListenerEventName(tsNode)
      }
      return ts.visitEachChild(tsNode, visit, _transformContext)
    }

    return tsSourceFile => {
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
const PREFIX = 'body:'
function normalize(eventName: string): string {
  if (!isBearerEvent(eventName)) {
    return eventName
  }
  if (eventName.startsWith(PREFIX)) {
    return `${PREFIX}${normalize(eventName.split(PREFIX)[1])}`
  }
  return Case.kebab(eventName)
}
