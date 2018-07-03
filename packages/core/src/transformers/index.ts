import * as ts from 'typescript'
import * as d from './declarations'
import gatherMetadata from './gather-metadata'
// import scopeTagNames from './scoped-tag-names'
import decoratorTransformers from './decorator-transformer'

export function Transformers(
  options: d.PluginOptions
): ts.TransformerFactory<ts.SourceFile>[] {
  const metadatas = { components: [] }

  return [
    gatherMetadata(metadatas),
    decoratorTransformers(options)
  ]
}

export default Transformers
