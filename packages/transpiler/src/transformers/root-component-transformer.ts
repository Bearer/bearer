import * as ts from 'typescript'
import * as Case from 'case'
import { Decorators } from './constants'
import { hasDecoratorNamed, decoratorNamed } from './decorator-helpers'

type TransformerOptions = {
  verbose?: true
}

export default function RootComponentTransformer({  }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    function visit(node: ts.Node): ts.Node {
      if (ts.isClassDeclaration(node) && hasDecoratorNamed(node, Decorators.RootComponent)) {
        let group = ''
        let name = ''
        const decorators = node.decorators.map(decorator => {
          if (decoratorNamed(decorator, Decorators.RootComponent)) {
            const obj = (decorator.expression as ts.CallExpression).arguments[0] as ts.ObjectLiteralExpression
            obj.properties.forEach(property => {
              if (
                property
                  .getChildAt(0)
                  .getFullText()
                  .trim() === 'group'
              ) {
                group = property
                  .getChildAt(2)
                  .getFullText()
                  .trim()
                  .slice(1, -1)
              }
              if (
                property
                  .getChildAt(0)
                  .getFullText()
                  .trim() === 'name'
              ) {
                name = property
                  .getChildAt(2)
                  .getFullText()
                  .trim()
                  .slice(1, -1)
              }
            })
            const tagComponent = [Case.kebab(group), name].join('-')
            const cssFileName = Case.pascal(group)
            return ts.updateDecorator(
              decorator,
              ts.createCall(ts.createIdentifier(Decorators.Component), undefined, [
                ts.createObjectLiteral(
                  [
                    ts.createPropertyAssignment('tag', ts.createStringLiteral(tagComponent)),
                    ts.createPropertyAssignment('styleUrl', ts.createStringLiteral(cssFileName + '.css')),
                    ts.createPropertyAssignment('shadow', ts.createTrue())
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
