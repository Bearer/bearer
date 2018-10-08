/*
 * Event Name scoping. make sure all events are scoped
 */
import * as ts from 'typescript'

import { Decorators, Env } from '../constants'
import { decoratorNamed, hasDecoratorNamed } from '../helpers/decorator-helpers'
type TransformerOptions = {
  verbose?: true
}

function updateEventDecorator(tsProperty: ts.PropertyDeclaration, scope: string): ts.PropertyDeclaration {
  // const expression: ts.CallExpression = tsDecorator.expression as ts.CallExpression
  const decorator = ts.createDecorator(
    ts.createCall(ts.createIdentifier(Decorators.Event), undefined, [
      ts.createObjectLiteral([
        ts.createPropertyAssignment(
          'eventName',
          ts.createLiteral(['bearer', process.env[Env.BEARER_SCENARIO_ID], scope, tsProperty.name.getText()].join(':'))
        )
      ])
    ])
  )

  return ts.updateProperty(
    tsProperty,
    [...tsProperty.decorators.map(deco => (decoratorNamed(deco, Decorators.Event) ? decorator : deco))],
    tsProperty.modifiers,
    tsProperty.name.getText(),
    tsProperty.questionToken || tsProperty.exclamationToken,
    tsProperty.type,
    tsProperty.initializer
  )
}

export default function EventNameScoping(_params: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    return tsSourceFile => {
      const groupName: string = 'group'

      function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
        if (ts.isPropertyDeclaration(tsNode) && hasDecoratorNamed(tsNode, Decorators.Event)) {
          return updateEventDecorator(tsNode, groupName)
        }
        return ts.visitEachChild(tsNode, visit, _transformContext)
      }

      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
