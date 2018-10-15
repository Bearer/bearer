import * as Case from 'case'
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { getDecoratorNamed, getExpressionFromDecorator, hasDecoratorNamed } from '../helpers/decorator-helpers'
import { TransformerOptions } from '../types'

export default function GatherMetadata({ metadata }: TransformerOptions): ts.TransformerFactory<ts.SourceFile> {
  function getTagNames(tagName: string): { initialTagName: string; finalTagName: string } {
    const finalTag =
      metadata.prefix && metadata.suffix
        ? [Case.kebab(metadata.prefix), Case.kebab(metadata.suffix), tagName].join('-')
        : tagName
    return {
      initialTagName: tagName,
      finalTagName: finalTag
    }
  }

  return _transformContext => {
    return tsSourceFile => {
      function visit(node: ts.Node): ts.Node {
        // Found Component
        if (ts.isClassDeclaration(node)) {
          if (hasDecoratorNamed(node, Decorators.Component)) {
            const component = getDecoratorNamed(node, Decorators.Component)
            const tag = getExpressionFromDecorator<ts.StringLiteral>(component, 'tag')
            metadata.registerComponent({
              fileName: tsSourceFile.fileName,
              classname: node.name.text,
              isRoot: false,
              ...getTagNames(tag.text)
            })
          }
          // Found RootComponent
          else if (hasDecoratorNamed(node, Decorators.RootComponent)) {
            const component = getDecoratorNamed(node, Decorators.RootComponent)
            const nameExpression = getExpressionFromDecorator<ts.StringLiteral>(component, 'role')
            const name = nameExpression ? nameExpression.text : ''
            const groupExpression = getExpressionFromDecorator<ts.StringLiteral>(component, 'group')
            const group = groupExpression ? groupExpression.text : ''
            const tag = [Case.kebab(group), name].join('-')
            metadata.registerComponent({
              fileName: tsSourceFile.fileName,
              classname: node.name.text,
              isRoot: true,
              ...getTagNames(tag),
              group
            })
          }
        }
        return node
      }
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
