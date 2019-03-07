/*
 * Checks if class is decorated with @Component decorator
 * and injects the `@Prop() INTEGRATION_ID: string;` into class definition
 *
 */
import * as ts from 'typescript'

import { Decorators, Env } from '../constants'
import { hasDecoratorNamed } from '../helpers/decorator-helpers'
import { TransformerOptions } from '../types'

import bearer from './bearer'

export function shouldProcessFile(tsSourceFile: ts.SourceFile): boolean {
  if (tsSourceFile.isDeclarationFile) {
    return false
  }
  return ts.forEachChild<boolean>(tsSourceFile, shouldInject)
}

export function shouldInject(node: ts.Node): boolean {
  return ts.isClassDeclaration(node) && hasDecoratorNamed(node, Decorators.Component)
}

export function retrieveIntegrationId(): string {
  return process.env[Env.BEARER_INTEGRATION_ID]
}

export default function componentTransformer({  }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return transformContext => {
    const integrationId = retrieveIntegrationId()

    function visit(node: ts.Node): ts.VisitResult<ts.Node> {
      // TODO: filter components which really need it
      if (shouldInject(node)) {
        return ts.visitEachChild(
          bearer.addBearerIntegrationIdAccessor(node as ts.ClassDeclaration, integrationId),
          visit,
          transformContext
        )
      }
      return ts.visitEachChild(node, visit, transformContext)
    }

    if (!integrationId) {
      console.warn('[BEARER]', 'No integration ID provided. Skipping integration ID injection')
    }

    return tsSourceFile => {
      if (!integrationId || !shouldProcessFile(tsSourceFile)) {
        return tsSourceFile
      }
      return visit(tsSourceFile) as ts.SourceFile
    }
  }
}
