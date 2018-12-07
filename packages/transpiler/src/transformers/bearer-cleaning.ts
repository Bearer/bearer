/*
 * Transformer boilerplate
 */
import * as ts from 'typescript'

import { Decorators, Types } from '../constants'
import { hasDecoratorNamed } from '../helpers/decorator-helpers'
import { getNodeName } from '../helpers/node-helpers'
import { TransformerOptions } from '../types'

const filteredImports = new Set([Decorators.Input, Decorators.Output, Decorators.RootComponent].map(v => v.toString()))

export default function bearerCleaning(_options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      // returning null will remove the node
      switch (tsNode.kind) {
        case ts.SyntaxKind.PropertyDeclaration: {
          const prop = tsNode as ts.PropertyDeclaration
          if (hasDecoratorNamed(prop, Decorators.Input) || hasDecoratorNamed(prop, Decorators.Output)) {
            return null
          }
          if (hasDecoratorNamed(prop, Decorators.Event)) {
            return removeEventEmitterType(prop)
          }
          break
        }

        case ts.SyntaxKind.ImportSpecifier: {
          const importSpecifier = tsNode as ts.ImportSpecifier
          if (filteredImports.has(getNodeName(importSpecifier))) {
            return null
          }
          break
        }
        default:
      }
      return ts.visitEachChild(tsNode, visit, _transformContext)
    }

    return tsSourceFile => {
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}

function removeEventEmitterType(property: ts.PropertyDeclaration): ts.PropertyDeclaration {
  if (!property.type) {
    return property
  }
  return ts.updateProperty(
    property,
    property.decorators,
    property.modifiers,
    property.name,
    property.questionToken,
    ts.createTypeReferenceNode(Types.EventEmitter, undefined),
    property.initializer
  )
}
