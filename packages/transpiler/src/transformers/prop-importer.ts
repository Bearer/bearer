/*
 * Checks if there is an `import { Component } from '@bearer/core'` in the source code
 * and adds the following if it is included:
 * 
 * `import { Prop } from '@bearer/core';`
 * 
 */
import * as ts from 'typescript'

import bearer from './bearer'

type TransformerOptions = {
  verbose?: true
}
export default function PropImporter({
  verbose
}: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  function log(...args) {
    if (verbose) {
      console.log.apply(this, args)
    }
  }

  return _transformContext => {
    return tsSourceFile => {
      if (
        bearer.hasImport(tsSourceFile, 'Component') &&
        !bearer.hasImport(tsSourceFile, 'Prop')
      ) {
        return ts.updateSourceFileNode(tsSourceFile, [
          ts.createImportDeclaration(
            undefined,
            undefined,
            ts.createImportClause(
              undefined,
              ts.createNamedImports([
                ts.createImportSpecifier(undefined, ts.createIdentifier('Prop'))
              ])
            ),
            ts.createLiteral('@bearer/core')
          ),
          ...tsSourceFile.statements
        ])
      }
      return tsSourceFile
    }
  }
}
