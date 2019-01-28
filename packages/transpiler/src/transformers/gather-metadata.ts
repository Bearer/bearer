import * as Case from 'case'
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { getDecoratorNamed, getExpressionFromDecorator, hasDecoratorNamed } from '../helpers/decorator-helpers'
import { TransformerOptions } from '../types'
import debug from '../logger'
const logger = debug.extend('gather-metadata')

export default function gatherMetadata({ metadata }: TransformerOptions): ts.TransformerFactory<ts.SourceFile> {
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
            const names = getTagNames(tag.text)
            metadata.registerComponent({
              fileName: tsSourceFile.fileName,
              classname: node.name.text,
              isRoot: false,
              ...names
            })
            logger('Registered %s: new tag name => ', names.initialTagName, names.finalTagName)
          }
          // Found RootComponent
          else if (hasDecoratorNamed(node, Decorators.RootComponent)) {
            const component = getDecoratorNamed(node, Decorators.RootComponent)
            const nameExpression = getExpressionFromDecorator<ts.StringLiteral>(component, 'role')
            const name = nameExpression ? nameExpression.text : ''
            const groupExpression = getExpressionFromDecorator<ts.StringLiteral>(component, 'group')
            const group = groupExpression ? groupExpression.text : ''
            const tag = [Case.kebab(group), name].join('-')
            const names = getTagNames(tag)
            metadata.registerComponent({
              fileName: tsSourceFile.fileName,
              classname: node.name.text,
              isRoot: true,
              ...names,
              group
            })
            logger('Registered RootComponent %s: new tag name => ', names.initialTagName, names.finalTagName)
          }
        }
        return node
      }
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
