/*
 * @Component() 
 * class StarWarsMovies {}
 * 
 * becomes
 * 
 * @Component()
 * class StarWarsMovies {
 * 
 *  @Prop({ context: 'bearer' }) bearerContext: string;
 *  @Prop() setupId: string;
 * 
 *  componentDidLoad() {
 *    if(this.setupId) {
 *      this.bearerContext.setupId = this.setupId
 *    }
 *  }
 * }
 * 
 */
import * as ts from 'typescript'

import decorator from './decorator-helpers'
import bearer from './bearer'

type TransformerOptions = {
  verbose?: true
}

function injectContext(node: ts.ClassDeclaration): ts.Node {
  const withContextProp = bearer.addBearerContextProp(node)
  const withSetupProp = bearer.addSetupIdProp(withContextProp)
  return bearer.addComponentDidLoad(withSetupProp)
}

export default function ComponentTransformer({
  verbose
}: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return transformContext => {
    function visit(node: ts.Node): ts.VisitResult<ts.Node> {
      if (
        ts.isClassDeclaration(node) &&
        decorator.classDecoratedWithName(
          node as ts.ClassDeclaration,
          'Component'
        )
      ) {
        return ts.visitEachChild(
          injectContext(node as ts.ClassDeclaration),
          visit,
          transformContext
        )
      }
      return ts.visitEachChild(node, visit, transformContext)
    }

    return tsSourceFile => {
      return visit(tsSourceFile) as ts.SourceFile
    }
  }
}
