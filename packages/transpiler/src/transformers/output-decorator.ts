/*
 * Transformer boilerplate
 */
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { hasDecoratorNamed } from '../helpers/decorator-helpers'
import { getNodeName } from '../helpers/node-helpers'
import { TransformerOptions } from '../types'

export default function OutputDecorator(_options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      return ts.visitEachChild(tsNode, visit, _transformContext)
    }

    return tsSourceFile => {
      if (tsSourceFile.isDeclarationFile) {
        return tsSourceFile
      }

      const outputsMeta = retrieveOutputsMetas(tsSourceFile)

      if (!outputsMeta.length) {
        return tsSourceFile
      }

      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }

    function retrieveOutputsMetas(tsSourceFile: ts.SourceFile): Array<OutputMeta> {
      const outputs: Array<OutputMeta> = []

      const visitor = (tsNode: ts.Node) => {
        if (ts.isPropertyDeclaration(tsNode) && hasDecoratorNamed(tsNode, Decorators.Input)) {
          const name = getNodeName(tsNode)
          outputs.push({
            emitMethodName: outputEventName(name), // TODO: retrieve from options
            propDeclarationName: name,
            typeIdentifier: tsNode.type,
            watchedPropName: name // TODO: retrieve from options
          })
        }
        return ts.visitEachChild(tsNode, visitor, _transformContext)
      }

      ts.visitEachChild(tsSourceFile, visitor, _transformContext)

      return outputs
    }
  }
}

export function outputEventName(prefix: string, suffix?: string): string {
  const _suffix = suffix || 'Saved'
  return `${prefix}${_suffix.charAt(0).toUpperCase() + _suffix.slice(1)}`
}

type OutputMeta = {
  emitMethodName: string
  propDeclarationName: string
  typeIdentifier?: ts.TypeNode
  watchedPropName: string
}
