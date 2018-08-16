import * as Case from 'case'
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { getDecoratorNamed, getExpressionFromDecorator, hasDecoratorNamed } from '../helpers/decorator-helpers'
import { TransformerOptions } from '../types'

export default function GatherMetadata({ metadata }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    function visit(node: ts.Node): ts.Node {
      // Found Component
      if (ts.isClassDeclaration(node) && hasDecoratorNamed(node, Decorators.Component)) {
        const component = getDecoratorNamed(node, Decorators.Component)
        const tag = getExpressionFromDecorator<ts.StringLiteral>(component, 'tag')
        metadata.components.push({
          classname: node.name.text,
          isRoot: false,
          initialTagName: tag.text,
          finalTagName: tag.text
        })
        return node
      }

      // Found RootComponent
      if (ts.isClassDeclaration(node) && hasDecoratorNamed(node, Decorators.RootComponent)) {
        const component = getDecoratorNamed(node, Decorators.RootComponent)
        const nameExpression = getExpressionFromDecorator<ts.StringLiteral>(component, 'name')
        const name = nameExpression ? nameExpression.text : ''
        const groupExpression = getExpressionFromDecorator<ts.StringLiteral>(component, 'group')
        const group = groupExpression ? groupExpression.text : ''
        const tag = [Case.kebab(group), name].join('-')
        const finalTag =
          metadata.prefix && metadata.suffix
            ? [Case.kebab(metadata.prefix), tag, Case.kebab(metadata.suffix)].join('-')
            : tag
        metadata.components.push({
          classname: node.name.text,
          isRoot: true,
          initialTagName: tag,
          finalTagName: finalTag,
          group
        })
        return node
      }
      return node
    }
    return tsSourceFile => {
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
