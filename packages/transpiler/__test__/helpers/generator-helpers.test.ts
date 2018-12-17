import * as ts from 'typescript'

import { createFetcher } from '../../src/helpers/generator-helpers'
import { InputMeta } from '../../src/types'

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

const meta: InputMeta = {
  propDeclarationName: 'propDeclarationName',
  group: 'group',
  propertyReferenceIdName: 'propertyReferenceIdName',
  eventName: 'eventName',
  intentReferenceIdKeyName: 'intentRefernceIdKeyName',
  intentName: 'intentName',
  autoLoad: false,
  loadMethodName: 'loadMethodName',
  intentMethodName: 'intentMethodName',
  watcherName: 'watcherName'
}

it('generates stuff properly', () => {
  const sourceFile = ts.createSourceFile('tmp.ts', code, ts.ScriptTarget.Latest)
  const transformed = ts.transform(sourceFile, [transform])

  const printer = ts.createPrinter(
    {
      newLine: ts.NewLineKind.LineFeed
    },
    {
      onEmitNode: transformed.emitNodeWithNotification,
      substituteNode: transformed.substituteNode
    }
  )
  const result = printer.printBundle(ts.createBundle(transformed.transformed))
  transformed.dispose()

  expect(result).toMatchSnapshot()
})
