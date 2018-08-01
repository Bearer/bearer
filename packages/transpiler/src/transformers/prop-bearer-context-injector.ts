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

import { hasDecoratorNamed } from './decorator-helpers'
import bearer from './bearer'
import { Decorators } from './constants'

type TransformerOptions = {
  verbose?: true
}

export default function ComponentTransformer({  }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return transformContext => {
    function visit(node: ts.Node): ts.VisitResult<ts.Node> {
      if (ts.isClassDeclaration(node) && hasDecoratorNamed(node, Decorators.Component)) {
        return ts.visitEachChild(injectContext(node as ts.ClassDeclaration), visit, transformContext)
      }
      return ts.visitEachChild(node, visit, transformContext)
    }

    return tsSourceFile => {
      if (hasComponentDecorator(tsSourceFile)) {
        return visit(bearer.ensurePropImported(tsSourceFile)) as ts.SourceFile
      }
      return tsSourceFile
    }
  }
}

function injectContext(node: ts.ClassDeclaration): ts.Node {
  const withContextProp = bearer.ensureBearerContextInjected(node)
  const withSetupProp = bearer.addSetupIdProp(withContextProp)
  return bearer.addComponentDidLoad(withSetupProp)
}

function hasComponentDecorator(sourceFile: ts.SourceFile): boolean {
  if (sourceFile.isDeclarationFile) {
    return false
  }

  return ts.forEachChild(
    sourceFile,
    node => ts.isClassDeclaration(node) && hasDecoratorNamed(node, Decorators.Component)
  )
}
