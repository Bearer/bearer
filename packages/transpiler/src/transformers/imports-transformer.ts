import * as ts from 'typescript'

import { ensureComponentImported, ensureRootComponentNotImported, hasImport } from './bearer'
import { Decorators } from './constants'
type TransformerOptions = {
  verbose?: true
}
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
