/*
 * Checks if class is decorated with @Component decorator
 * and injects the `@Prop() SCENARIO_ID: string;` into class definition
 * 
 */
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { hasDecoratorNamed } from '../helpers/decorator-helpers'
import { TransformerOptions } from '../types'

import bearer from './bearer'

export default function ComponentTransformer({

}: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return transformContext => {
    const scenarioId = process.env.BEARER_SCENARIO_ID

    function visit(node: ts.Node): ts.VisitResult<ts.Node> {
      // TODO: filter components which really need it
      if (
        ts.isClassDeclaration(node) &&
        hasDecoratorNamed(node, Decorators.Component)
      ) {
        return ts.visitEachChild(
          bearer.addBearerScenarioIdAccessor(
            node as ts.ClassDeclaration,
            scenarioId
          ),
          visit,
          transformContext
        )
      }
      return ts.visitEachChild(node, visit, transformContext)
    }

    if (!scenarioId) {
      console.warn(
        '[BEARER]',
        'No scenario ID provided. Skipping scenario ID injection'
      )
    }

    return tsSourceFile => {
      if (!scenarioId) {
        return tsSourceFile
      }
      return visit(tsSourceFile) as ts.SourceFile
    }
  }
}
