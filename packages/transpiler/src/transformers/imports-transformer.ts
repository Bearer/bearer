import * as ts from 'typescript'

import { Decorators } from '../constants'
import { TransformerOptions } from '../types'

import { ensureComponentImported, ensureRootComponentNotImported, hasImport } from './bearer'

export default function ImportsImporter({  }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    return tsSourceFile => {
      if (hasImport(tsSourceFile, Decorators.RootComponent)) {
        return ensureRootComponentNotImported(ensureComponentImported(tsSourceFile))
      }
      return tsSourceFile
    }
  }
}
