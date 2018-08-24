import * as Case from 'case'
import * as ts from 'typescript'

import { Decorators } from '../constants'
import {
  getDecoratorNamed,
  getExpressionFromDecorator,
  hasDecoratorNamed
} from '../helpers/decorator-helpers'
import { FileTransformerOptions } from '../types'

import generateManifestFile from './generate-manifest-file'

export default function GatherMetadata(
  { metadata, outDir, srcDir }: FileTransformerOptions = { outDir: null }
): ts.TransformerFactory<ts.SourceFile> {
  function getTagNames(
    tagName: string
  ): { initialTagName: string; finalTagName: string } {
    const finalTag =
      metadata.prefix && metadata.suffix
        ? [
            Case.kebab(metadata.prefix),
            tagName,
            Case.kebab(metadata.suffix)
          ].join('-')
        : tagName
    return {
      initialTagName: tagName,
      finalTagName: finalTag
    }
  }

  return _transformContext => {
    function visit(node: ts.Node): ts.Node {
      // Found Component
      if (
        ts.isClassDeclaration(node) &&
        hasDecoratorNamed(node, Decorators.Component)
      ) {
        const component = getDecoratorNamed(node, Decorators.Component)
        const tag = getExpressionFromDecorator<ts.StringLiteral>(
          component,
          'tag'
        )
        metadata.components.push({
          classname: node.name.text,
          isRoot: false,
          ...getTagNames(tag.text)
        })
        return node
      }

      // Found RootComponent
      if (
        ts.isClassDeclaration(node) &&
        hasDecoratorNamed(node, Decorators.RootComponent)
      ) {
        const component = getDecoratorNamed(node, Decorators.RootComponent)
        const nameExpression = getExpressionFromDecorator<ts.StringLiteral>(
          component,
          'role'
        )
        const name = nameExpression ? nameExpression.text : ''
        const groupExpression = getExpressionFromDecorator<ts.StringLiteral>(
          component,
          'group'
        )
        const group = groupExpression ? groupExpression.text : ''
        const tag = [Case.kebab(group), name].join('-')

        metadata.components.push({
          classname: node.name.text,
          isRoot: true,
          ...getTagNames(tag),
          group
        })

        generateManifestFile({ metadata, outDir, srcDir })
        return node
      }
      return node
    }
    return tsSourceFile => {
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
