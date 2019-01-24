/*
 * Event Name scoping. make sure all events are scoped
 */
import * as ts from 'typescript'

import { BEARER, Decorators, Env, Properties } from '../constants'
import { decoratorNamed, hasDecoratorNamed } from '../helpers/decorator-helpers'
import { getNodeName } from '../helpers/node-helpers'
import { TransformerOptions } from '../types'

const SEPARATOR = '-'

export function prefixEvent(eventName: string): string {
  return [BEARER, process.env[Env.BEARER_SCENARIO_ID], eventName].join(SEPARATOR)
}

export function eventName(name: string, scope = 'no-group') {
  return prefixEvent([scope, name].filter(el => el && el.trim()).join(SEPARATOR))
}

function updateEventDecorator(tsProperty: ts.PropertyDeclaration, scope: string): ts.PropertyDeclaration {
  const decorator = ts.createDecorator(
    ts.createCall(ts.createIdentifier(Decorators.Event), undefined, [
      ts.createObjectLiteral([
        ts.createPropertyAssignment(Properties.eventName, ts.createLiteral(eventName(getNodeName(tsProperty), scope)))
      ])
    ])
  )

  return ts.updateProperty(
    tsProperty,
    [...tsProperty.decorators.map(deco => (decoratorNamed(deco, Decorators.Event) ? decorator : deco))],
    tsProperty.modifiers,
    getNodeName(tsProperty),
    tsProperty.questionToken || tsProperty.exclamationToken,
    tsProperty.type,
    tsProperty.initializer
  )
}

function updatedListenDecoratorOrDecorator(tsDecorator: ts.Decorator): ts.Decorator {
  if (!decoratorNamed(tsDecorator, Decorators.Listen)) {
    return tsDecorator
  }

  const listenedEvent: ts.StringLiteral = (tsDecorator.expression as ts.CallExpression).arguments[0] as ts.StringLiteral

  // possible values
  //    group:eventName
  //    body:group:eventName
  let [body, group, name] = listenedEvent.text.toString().split(':')
  if (body !== 'body') {
    name = group
    group = body
    body = null
  }
  let scopedName = eventName(name, group)
  if (body) {
    scopedName = ['body', scopedName].join(':')
  }
  return ts.createDecorator(
    ts.createCall(ts.createIdentifier(Decorators.Listen), undefined, [ts.createLiteral(scopedName)])
  )
}

export default function eventNameScoping({ metadata }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    return tsSourceFile => {
      const meta = metadata.findComponentFrom(tsSourceFile)
      const groupName = (meta && meta.group) || 'global'

      function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
        if (ts.isPropertyDeclaration(tsNode) && hasDecoratorNamed(tsNode, Decorators.Event)) {
          return updateEventDecorator(tsNode, groupName)
        }
        if (ts.isMethodDeclaration(tsNode) && hasDecoratorNamed(tsNode, Decorators.Listen)) {
          return ts.updateMethod(
            tsNode,
            [...tsNode.decorators.map(updatedListenDecoratorOrDecorator)],
            tsNode.modifiers,
            tsNode.asteriskToken,
            tsNode.name,
            tsNode.questionToken,
            tsNode.typeParameters,
            tsNode.parameters,
            tsNode.type,
            tsNode.body
          )
        }
        return ts.visitEachChild(tsNode, visit, _transformContext)
      }

      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
