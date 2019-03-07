import * as ts from 'typescript'

import { Decorators } from '../constants'
import { decoratorNamed, getExpressionFromDecorator, hasDecoratorNamed } from '../helpers/decorator-helpers'
import { TransformerOptions } from '../types'

import { ensureImportsFromCore, hasImport } from './bearer'
import debug from '../logger'

const logger = debug.extend('root-component')

export default function rootComponentTransformer({ metadata }: TransformerOptions = {}): ts.TransformerFactory<
  ts.SourceFile
> {
  return _transformContext => {
    function visit(node: ts.Node): ts.Node {
      if (ts.isClassDeclaration(node) && hasDecoratorNamed(node, Decorators.RootComponent)) {
        const decorators = node.decorators.map(decorator => {
          if (decoratorNamed(decorator, Decorators.RootComponent)) {
            const metadatum = metadata.components.find(component => component.classname === node.name.text)
            const shadowExp = getExpressionFromDecorator<ts.BooleanLiteral>(decorator, 'shadow')
            const styleUrl = getExpressionFromDecorator<ts.StringLiteral>(decorator, 'styleUrl')
            const options = [
              ts.createPropertyAssignment('tag', ts.createStringLiteral(metadatum.finalTagName)),
              ts.createPropertyAssignment('shadow', shadowExp || ts.createTrue())
            ]
            if (styleUrl) {
              options.push(ts.createPropertyAssignment('styleUrl', ts.createStringLiteral(styleUrl.text)))
            }

            return ts.updateDecorator(
              decorator,
              ts.createCall(ts.createIdentifier(Decorators.Component), undefined, [
                ts.createObjectLiteral(options, true)
              ])
            )
          }
          return decorator
        })

        return ts.updateClassDeclaration(
          node,
          decorators,
          node.modifiers,
          node.name,
          node.typeParameters,
          node.heritageClauses,
          node.members
        )
      }
      return node
    }
    return tsSourceFile => {
      if (tsSourceFile.isDeclarationFile || !hasImport(tsSourceFile, Decorators.RootComponent)) {
        return tsSourceFile
      }
      logger('processing %s', tsSourceFile.fileName)
      return ts.visitEachChild(ensureImportsFromCore(tsSourceFile, [Decorators.Component]), visit, _transformContext)
    }
  }
}
