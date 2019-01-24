import * as ts from 'typescript'

import { Decorators, Properties } from '../constants'
import { decoratorNamed } from '../helpers/decorator-helpers'
import * as Case from 'case'

export default function eventNameNormalizer(): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    function normalizeEventName(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      if (
        ts.isPropertyAssignment(tsNode) &&
        (tsNode.name as ts.Identifier).escapedText.toString() === Properties.eventName
      ) {
        const normalizedName = Case.kebab((tsNode.initializer as ts.StringLiteral).text.toString())
        return ts.createPropertyAssignment(tsNode.name, ts.createLiteral(normalizedName))
      }
      return ts.visitEachChild(tsNode, normalizeEventName, _transformContext)
    }

    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      if (ts.isDecorator(tsNode) && decoratorNamed(tsNode, Decorators.Event)) {
        return normalizeEventName(tsNode)
      }
      return ts.visitEachChild(tsNode, visit, _transformContext)
    }

    return tsSourceFile => {
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
