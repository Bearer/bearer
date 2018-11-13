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

import { Decorators } from '../constants'
import { hasDecoratorNamed } from '../helpers/decorator-helpers'
import { TransformerOptions } from '../types'

import bearer, { ensureImportsFromCore } from './bearer'

export default function ComponentTransformer(_options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return transformContext => {
    function visit(node: ts.Node): ts.VisitResult<ts.Node> {
      if (ts.isClassDeclaration(node) && hasDecoratorNamed(node, Decorators.Component)) {
        return ts.visitEachChild(injectContext(node as ts.ClassDeclaration), visit, transformContext)
      }
      return ts.visitEachChild(node, visit, transformContext)
    }

    return tsSourceFile => {
      if (hasComponentDecorator(tsSourceFile)) {
        return visit(ensureImportsFromCore(tsSourceFile, [Decorators.Prop])) as ts.SourceFile
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
