/*
 * Transformer boilerplate
 */
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { TransformerOptions } from '../types'

import { ensureHasNotImportFromCore, hasImport } from './bearer'

export default function bearerCleaning(_options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      return ts.visitEachChild(tsNode, visit, _transformContext)
    }

    return tsSourceFile => {
      const cleanedSourceFile = removeRootComponentDecorator(tsSourceFile)
      return ts.visitEachChild(cleanedSourceFile, visit, _transformContext)
    }
  }
}

function removeRootComponentDecorator(tsSourceFile: ts.SourceFile) {
  if (hasImport(tsSourceFile, Decorators.RootComponent)) {
    return ensureRootComponentNotImported(tsSourceFile)
  }
  return tsSourceFile
}

export function ensureRootComponentNotImported(tsSourceFile: ts.SourceFile): ts.SourceFile {
  return ensureHasNotImportFromCore(tsSourceFile, Decorators.RootComponent)
}
