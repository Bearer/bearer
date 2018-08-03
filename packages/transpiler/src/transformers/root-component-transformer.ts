import * as ts from 'typescript'
import * as Case from 'case'
import { Decorators } from '../constants'
import { hasDecoratorNamed, decoratorNamed, getExpressionFromDecorator } from '../helpers/decorator-helpers'
import { TransformerOptions } from '../types'

export default function RootComponentTransformer({ metadata }: TransformerOptions = {}): ts.TransformerFactory<
  ts.SourceFile
> {
  return _transformContext => {
    function visit(node: ts.Node): ts.Node {
      if (ts.isClassDeclaration(node) && hasDecoratorNamed(node, Decorators.RootComponent)) {
        const decorators = node.decorators.map(decorator => {
          if (decoratorNamed(decorator, Decorators.RootComponent)) {
            const metadatum = metadata.components.find(component => component.classname === node.name.text)
            const cssFileName = Case.pascal(metadatum.group)
            const shadowExp = getExpressionFromDecorator<ts.BooleanLiteral>(decorator, 'shadow')
            return ts.updateDecorator(
              decorator,
              ts.createCall(ts.createIdentifier(Decorators.Component), undefined, [
                ts.createObjectLiteral(
                  [
                    ts.createPropertyAssignment('tag', ts.createStringLiteral(metadatum.initialTagName)),
                    ts.createPropertyAssignment('styleUrl', ts.createStringLiteral(cssFileName + '.css')),
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
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
