import * as ts from 'typescript'

import { addAutoLoad, createFetcher, createLoadResourceMethod } from '../../src/helpers/generator-helpers'
import { CreateFetcherMeta } from '../../src/types'
import { runTransformers } from '../utils/helpers'

function dummyTransformer(fun: (meta: any) => ts.Node, meta: any) {
  return (context: ts.TransformationContext) => {
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
          [...prop.members, fun(meta) as ts.ClassElement]
        )
      }

      return ts.visitEachChild(node, updateClass, context)
    }

    return (sourceFile: ts.SourceFile): ts.SourceFile => {
      return ts.visitNode(sourceFile, updateClass)
    }
  }
}

describe('createFetcher', () => {
  it('generates stuff properly', () => {
    const code = `
     class C {}
    `

    const meta: CreateFetcherMeta = {
      intentName: 'intentName',
      intentMethodName: 'intentMethodName'
    }

    expect(runTransformers(code, [dummyTransformer(createFetcher, meta)])).toMatchSnapshot()
  })
})

describe('createLoadResourceMethod', () => {
  it('generates load resource method properly', () => {
    const code = `
     class C {}
    `

    const meta = {
      propertyReferenceIdName: 'propertyReferenceIdName',
      typeIdentifier: ts.createIdentifier('SomeType'),
      propDeclarationName: 'propDeclarationName',
      intentMethodName: 'intentMethodName',
      intentReferenceIdKeyName: 'intentReferenceIdKeyName',
      loadMethodName: 'loadMethodName'
    }
    expect(runTransformers(code, [dummyTransformer(createLoadResourceMethod, meta)])).toMatchSnapshot()
  })
})

describe('addAutoLoad', () => {
  it('generates class with autoloaded feature', () => {
    function dummyTransformer(fun: (node: ts.ClassDeclaration, meta: any) => ts.Node, meta: any) {
      return (context: ts.TransformationContext) => {
        const updateClass: ts.Visitor = (node: ts.Node) => {
          if (ts.isClassDeclaration(node)) {
            const prop = node as ts.ClassDeclaration
            return fun(prop, meta)
          }

          return ts.visitEachChild(node, updateClass, context)
        }

        return (sourceFile: ts.SourceFile): ts.SourceFile => {
          return ts.visitNode(sourceFile, updateClass)
        }
      }
    }
    const meta = {
      propertyReferenceIdName: 'propertyReferenceIdName',
      autoLoad: true,
      loadMethodName: 'loadMethodName'
    }
    expect(runTransformers('class C {}', [dummyTransformer(addAutoLoad, meta)])).toMatchSnapshot()
  })
})
