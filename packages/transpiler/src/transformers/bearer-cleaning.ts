/*
 * Transformer boilerplate
 */
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { TransformerOptions } from '../types'

import { ensureHasNotImportFromCore } from './bearer'

export default function bearerCleaning(_options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      return ts.visitEachChild(tsNode, visit, _transformContext)
    }

    return tsSourceFile => {
      const cleanedSourceFile = removeBearerDecorators(tsSourceFile)
      return ts.visitEachChild(cleanedSourceFile, visit, _transformContext)
    }
  }
}

function removeBearerDecorators(tsSourceFile: ts.SourceFile) {
  return [Decorators.RootComponent, Decorators.Input].reduce((sourceFile, importerdDecorator) => {
    return ensureHasNotImportFromCore(sourceFile, importerdDecorator)
  }, tsSourceFile)
}
