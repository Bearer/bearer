import * as ts from 'typescript'
import decoratorTransformers from './decorator-transformer'

export function Transformers(options): ts.TransformerFactory<ts.SourceFile>[] {
  return [decoratorTransformers(options)]
}

export default Transformers
