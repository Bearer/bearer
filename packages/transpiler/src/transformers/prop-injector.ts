import * as ts from 'typescript'

import decorator from './decorator-helpers'
import bearer from './bearer'
import { getSourceCode } from '../utils'

type TransformerOptions = {
  verbose?: true
}

export default function ComponentTransformer({
  verbose
}: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  function log(...args) {
    if (verbose) {
      console.log.apply(this, args)
    }
  }

  var propAdded = false
  return transformContext => {
    function visit(node: ts.Node): ts.VisitResult<ts.Node> {
      switch (node.kind) {
        case ts.SyntaxKind.ClassDeclaration: {
          propAdded = true
          if (
            decorator.classDecoratedWithName(
              node as ts.ClassDeclaration,
              'Component'
            )
          ) {
            return ts.visitEachChild(
              bearer.addBearerIdProp(node as ts.ClassDeclaration),
              visit,
              transformContext
            )
          } else {
            return ts.visitEachChild(node, visit, transformContext)
          }
        }
      }
      return ts.visitEachChild(node, visit, transformContext)
    }

    return tsSourceFile => {
      return visit(tsSourceFile) as ts.SourceFile
    }
  }
}
