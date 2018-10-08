/*
 * Event Name scoping. make sure all events are scoped
 */
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { hasDecoratorNamed } from '../helpers/decorator-helpers'
type TransformerOptions = {
  verbose?: true
}

function updateEventDecorator(tsProperty: ts.PropertyDeclaration): ts.PropertyDeclaration {
  // const expression: ts.CallExpression = tsDecorator.expression as ts.CallExpression
  return ts.updateProperty(
    tsProperty,
    tsProperty.decorators,
    tsProperty.modifiers,
    tsProperty.name.getText(),
    tsProperty.questionToken || tsProperty.exclamationToken,
    tsProperty.type,
    tsProperty.initializer
  )
}

export default function EventNameScoping(_params: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      if (ts.isPropertyDeclaration(tsNode) && hasDecoratorNamed(tsNode, Decorators.Event)) {
        return updateEventDecorator(tsNode)
      }
      return ts.visitEachChild(tsNode, visit, _transformContext)
    }

    return tsSourceFile => {
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
