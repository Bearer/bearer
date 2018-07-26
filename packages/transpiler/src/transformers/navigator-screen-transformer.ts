/*
 * Rewrite navigator-screen if they do not use renderFunc
 */
import * as ts from 'typescript'

type TransformerOptions = {
  verbose?: true
}
export default function PropImporter({ verbose }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    return tsSourceFile => {
      return tsSourceFile
    }
  }
}
