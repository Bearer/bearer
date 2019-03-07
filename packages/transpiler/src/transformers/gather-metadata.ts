import * as Case from 'case'
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { getDecoratorNamed, getExpressionFromDecorator, hasDecoratorNamed } from '../helpers/decorator-helpers'
import { TransformerOptions } from '../types'
import debug from '../logger'
const logger = debug.extend('gather-metadata')

export default function gatherMetadata({ metadata }: TransformerOptions): ts.TransformerFactory<ts.SourceFile> {
  function getTagNames(tagName: string): { initialTagName: string; finalTagName: string } {
    const parts = [Case.kebab(metadata.prefix) || 'bearer']
    if (metadata.suffix) {
      parts.push(Case.kebab(metadata.suffix))
    }
    parts.push(tagName)
    const finalTag = parts.join('-')
    return {
      initialTagName: tagName,
      finalTagName: finalTag
    }
  }

  return _transformContext => {
    return tsSourceFile => {
      logger('processing %s', tsSourceFile.fileName)
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
            const name = getExpressionFromDecorator<ts.StringLiteral>(component, 'name')

            const names = getTagNames(name.text.toString())
            metadata.registerComponent({
              fileName: tsSourceFile.fileName,
              classname: node.name.text,
              isRoot: true,
              name: name.text,
              ...names
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
