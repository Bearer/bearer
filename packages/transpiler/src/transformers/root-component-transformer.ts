import * as Case from 'case'
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { decoratorNamed, getExpressionFromDecorator, hasDecoratorNamed } from '../helpers/decorator-helpers'
import { TransformerOptions } from '../types'

import { ensureImportsFromCore, hasImport } from './bearer'

export default function rootComponentTransformer({ metadata }: TransformerOptions = {}): ts.TransformerFactory<
  ts.SourceFile
> {
  return _transformContext => {
    function visit(node: ts.Node): ts.Node {
      if (ts.isClassDeclaration(node) && hasDecoratorNamed(node, Decorators.RootComponent)) {
        const decorators = node.decorators.map(decorator => {
          if (decoratorNamed(decorator, Decorators.RootComponent)) {
            const metadatum = metadata.components.find(component => component.classname === node.name.text)
            const cssFileName = metadatum.group.charAt(0) + Case.camel(metadatum.group).substr(1)
            const shadowExp = getExpressionFromDecorator<ts.BooleanLiteral>(decorator, 'shadow')
            return ts.updateDecorator(
              decorator,
              ts.createCall(ts.createIdentifier(Decorators.Component), undefined, [
                ts.createObjectLiteral(
                  [
                    ts.createPropertyAssignment('tag', ts.createStringLiteral(metadatum.finalTagName)),
                    ts.createPropertyAssignment('styleUrl', ts.createStringLiteral([cssFileName, 'css'].join('.'))),
                    ts.createPropertyAssignment('shadow', shadowExp || ts.createTrue())
                  ],
                  true
                )
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
      return ts.visitEachChild(ensureImportsFromCore(tsSourceFile, [Decorators.Component]), visit, _transformContext)
    }
  }
}
