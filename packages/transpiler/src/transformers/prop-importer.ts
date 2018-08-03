/*
 * Checks if there is an `import { Component } from '@bearer/core'` in the source code
 * and adds the following if it is included:
 * 
 * `import { Prop } from '@bearer/core';`
 * 
 */
import * as ts from 'typescript'

import { ensurePropImported, hasImport } from './bearer'
import { Decorators } from '../constants'
import { TransformerOptions } from '../types'

export default function PropImporter({  }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    return tsSourceFile => {
      if (hasImport(tsSourceFile, Decorators.Component) && !hasImport(tsSourceFile, Decorators.Prop)) {
        return ensurePropImported(tsSourceFile)
      }
      return tsSourceFile
    }
  }
}
