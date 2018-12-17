import * as ts from 'typescript'

import { createFetcher } from '../../src/helpers/generator-helpers'
import { CreateFetcherMeta } from '../../src/types'
import { runTransformer } from '../utils/helpers'

function transform(context: ts.TransformationContext) {
  const updateClass: ts.Visitor = (node: ts.Node) => {
    if (ts.isClassDeclaration(node)) {
      const prop = node as ts.ClassDeclaration

      return ts.updateClassDeclaration(
        prop,
        prop.decorators,
        prop.modifiers,
        prop.name,
        prop.typeParameters,
        prop.heritageClauses,
        [...prop.members, createFetcher(meta)]
      )
    }

    return ts.visitEachChild(node, updateClass, context)
  }

  return (sourceFile: ts.SourceFile): ts.SourceFile => {
    return ts.visitNode(sourceFile, updateClass)
  }
}

const code = `
class C {}
`

const meta: CreateFetcherMeta = {
  intentName: 'intentName',
  intentMethodName: 'intentMethodName'
}

it('generates stuff properly', () => {
  expect(runTransformer(code, transform)).toMatchSnapshot()
})
