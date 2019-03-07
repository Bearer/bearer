/*
 * Event Name scoping. make sure all events are scoped
 */
import * as ts from 'typescript'

import { BEARER, Decorators, Env, Properties } from '../constants'
import { decoratorNamed, hasDecoratorNamed } from '../helpers/decorator-helpers'
import { getNodeName } from '../helpers/node-helpers'
import { TransformerOptions } from '../types'
import { normalize } from './event-name-normalizer'

const SEPARATOR = '-'
export const GLOBAL_EVENT_PREXIX = 'body:'

export function isBearerEvent(eventName: string): boolean {
  return new RegExp(`^(${GLOBAL_EVENT_PREXIX})?${BEARER}${SEPARATOR}`).test(eventName)
}

export function prefixEvent(eventName: string): string {
  return [BEARER, process.env[Env.BEARER_SCENARIO_ID], eventName].join(SEPARATOR)
}

export function eventName(name: string) {
  return normalize(prefixEvent(name))
}

function updateEventDecorator(tsProperty: ts.PropertyDeclaration): ts.PropertyDeclaration {
  const decorator = ts.createDecorator(
    ts.createCall(ts.createIdentifier(Decorators.Event), undefined, [
      ts.createObjectLiteral([
        ts.createPropertyAssignment(Properties.eventName, ts.createLiteral(eventName(getNodeName(tsProperty))))
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
  let [body, name] = listenedEvent.text.toString().split(':')
  if (body !== 'body') {
    name = body
    body = null
  }
  let scopedName = eventName(name)
  if (body) {
    scopedName = ['body', scopedName].join(':')
  }
  return ts.createDecorator(
    ts.createCall(ts.createIdentifier(Decorators.Listen), undefined, [ts.createLiteral(scopedName)])
  )
}

export default function eventNameScoping({  }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    return tsSourceFile => {
      function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
        if (ts.isPropertyDeclaration(tsNode) && hasDecoratorNamed(tsNode, Decorators.Event)) {
          return updateEventDecorator(tsNode)
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
