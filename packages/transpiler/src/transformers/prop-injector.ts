/*
 * Checks if class is decorated with @Component decorator
 * and injects the `@Prop() BEARER_ID: string;` into class definition
 * 
 */
import * as ts from 'typescript'

import { hasDecoratorNamed } from '../helpers/decorator-helpers'
import bearer from './bearer'
import { Decorators } from '../constants'
import { TransformerOptions } from '../types'

export default function ComponentTransformer({  }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return transformContext => {
    function visit(node: ts.Node): ts.VisitResult<ts.Node> {
      if (ts.isClassDeclaration(node) && hasDecoratorNamed(node, Decorators.Component)) {
        return ts.visitEachChild(bearer.addBearerIdProp(node as ts.ClassDeclaration), visit, transformContext)
      }
      return ts.visitEachChild(node, visit, transformContext)
    }

    return tsSourceFile => {
      return visit(tsSourceFile) as ts.SourceFile
    }
  }
}
