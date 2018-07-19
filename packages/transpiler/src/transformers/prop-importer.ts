/*
 * Checks if there is an `import { Component } from '@bearer/core'` in the source code
 * and adds the following if it is included:
 * 
 * `import { Prop } from '@bearer/core';`
 * 
 */
import * as ts from 'typescript'

import { ensurePropImported, hasImport } from './bearer'

type TransformerOptions = {
  verbose?: true
}
export default function PropImporter({
  verbose
}: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    return tsSourceFile => {
      if (
        hasImport(tsSourceFile, 'Component') &&
        !hasImport(tsSourceFile, 'Prop')
      ) {
        return ensurePropImported(tsSourceFile)
      }
      return tsSourceFile
    }
  }
}
