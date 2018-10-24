/*
 * Input Transformer
 */
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { hasDecoratorNamed } from '../helpers/decorator-helpers'
import { getNodeName } from '../helpers/node-helpers'
import { TransformerOptions } from '../types'

import { ensurePropImported } from './bearer'

export default function InputDecorator(_options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    return tsSourceFile => {
      if (tsSourceFile.isDeclarationFile) {
        return tsSourceFile
      }
      // gather inputs information
      const inputsMeta = retrieveInputsMetas(tsSourceFile)

      if (!inputsMeta.length) {
        return tsSourceFile
      }
      const withImports = ensurePropImported(tsSourceFile)
      return ts.visitEachChild(withImports, replaceInputVisitor(inputsMeta), _transformContext)
    }

    function retrieveInputsMetas(sourcefile: ts.SourceFile): Array<InputMeta> {
      const inputs: Array<InputMeta> = []

      const visitor = (tsNode: ts.Node) => {
        if (ts.isPropertyDeclaration(tsNode) && hasDecoratorNamed(tsNode, Decorators.Input)) {
          const name = getNodeName(tsNode)
          inputs.push({
            propDeclarationName: name,
            scope: 'string',
            propName: `${name}RefId`,
            eventName: `${name}:saved`,
            intentName: `get${name}`,
            autoUpdate: true,
            typeIdentifier: tsNode.type
          })
        }
        return ts.visitEachChild(tsNode, visitor, _transformContext)
      }

      ts.visitEachChild(sourcefile, visitor, _transformContext)

      return inputs
    }

    function injectInputStatements(tsClass: ts.ClassDeclaration, inputsMeta: Array<InputMeta>): ts.ClassDeclaration {
      const newMembers = inputsMeta.reduce(
        (members, meta) => {
          const inputMembers = [
            // create @State()
            ts.createProperty(
              [ts.createDecorator(ts.createCall(ts.createIdentifier(Decorators.State), undefined, undefined))],
              undefined,
              ts.createIdentifier(meta.propDeclarationName),
              undefined,
              meta.typeIdentifier,
              undefined
            )
          ]
          return members.concat(inputMembers)
        },
        [...tsClass.members]
      )

      return ts.updateClassDeclaration(
        tsClass,
        tsClass.decorators,
        tsClass.modifiers,
        tsClass.name,
        tsClass.typeParameters,
        tsClass.heritageClauses,
        newMembers
      )
    }

    function replaceInputVisitor(inputsMeta: Array<InputMeta>): (tsNode: ts.Node) => ts.VisitResult<ts.Node> {
      return (tsNode: ts.Node) => {
        // remove input usage
        if (ts.isPropertyDeclaration(tsNode) && hasDecoratorNamed(tsNode, Decorators.Input)) {
          return null
        }

        if (ts.isClassDeclaration(tsNode)) {
          return ts.visitEachChild(
            injectInputStatements(tsNode, inputsMeta),
            replaceInputVisitor(inputsMeta),
            _transformContext
          )
        }
        return ts.visitEachChild(tsNode, replaceInputVisitor(inputsMeta), _transformContext)
      }
    }
  }
}

type InputMeta = {
  propDeclarationName: string
  scope: string
  propName: string
  eventName: string
  intentName: string
  autoUpdate: boolean
  typeIdentifier?: ts.TypeNode
}
