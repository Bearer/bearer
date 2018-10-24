/*
 * Transformer boilerplate
 */
import * as ts from 'typescript'

type TransformerOptions = {
  verbose?: true
}

export default function CHANGEME(_options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      return ts.visitEachChild(tsNode, visit, _transformContext)
    }

    return tsSourceFile => {
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
