/*
 * Transformer boilerplate
 */
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { TransformerOptions } from '../types'

import { ensureNotImportedFromCore } from './bearer'

export default function bearerCleaning(_options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      return ts.visitEachChild(tsNode, visit, _transformContext)
    }

    return tsSourceFile => {
      const cleanedSourceFile = ensureNotImportedFromCore(tsSourceFile, [Decorators.RootComponent, Decorators.Input])
      return ts.visitEachChild(cleanedSourceFile, visit, _transformContext)
    }
  }
}
