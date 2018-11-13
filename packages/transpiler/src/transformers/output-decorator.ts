/*
 * Transformer boilerplate
 */
import * as ts from 'typescript'

import { Decorators, Types } from '../constants'
import { hasDecoratorNamed } from '../helpers/decorator-helpers'
import { getNodeName } from '../helpers/node-helpers'
import { TransformerOptions } from '../types'

export default function OutputDecorator(_options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    return tsSourceFile => {
      if (tsSourceFile.isDeclarationFile) {
        return tsSourceFile
      }

      const outputsMeta = retrieveOutputsMetas(tsSourceFile)

      if (!outputsMeta.length) {
        return tsSourceFile
      }

      return ts.visitEachChild(tsSourceFile, visit(outputsMeta), _transformContext)
    }

    function retrieveOutputsMetas(tsSourceFile: ts.SourceFile): Array<OutputMeta> {
      const outputs: Array<OutputMeta> = []

      const visitor = (tsNode: ts.Node) => {
        if (ts.isPropertyDeclaration(tsNode) && hasDecoratorNamed(tsNode, Decorators.Output)) {
          const name = getNodeName(tsNode)
          outputs.push({
            emitMethodName: outputEventName(name), // TODO: retrieve from options
            propDeclarationName: name,
            typeIdentifier: tsNode.type,
            watchedPropName: name // TODO: retrieve from options
          })
        }
        return ts.visitEachChild(tsNode, visitor, _transformContext)
      }

      ts.visitEachChild(tsSourceFile, visitor, _transformContext)

      return outputs
    }

    function visit(outputsMeta: Array<OutputMeta>) {
      return (tsNode: ts.Node): ts.VisitResult<ts.Node> => {
        if (ts.isClassDeclaration(tsNode)) {
          return injectOuputStatements(tsNode, outputsMeta)
        }
        return ts.visitEachChild(tsNode, visit(outputsMeta), _transformContext)
      }
    }

    function injectOuputStatements(tsClass: ts.ClassDeclaration, outputsMeta: Array<OutputMeta>): ts.ClassDeclaration {
      const newMembers = outputsMeta.reduce(
        (members, meta) => {
          const inputMembers = [createState(meta), createEvent(meta), createWatcher(meta)]
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

    function createEvent(meta: OutputMeta): ts.PropertyDeclaration {
      return ts.createProperty(
        [
          ts.createDecorator(
            ts.createCall(ts.createIdentifier(Decorators.Event), undefined, [
              ts.createStringLiteral(meta.watchedPropName)
            ])
          )
        ],
        undefined,
        meta.emitMethodName,
        undefined,
        ts.createTypeReferenceNode(ts.createIdentifier(Types.EventEmitter), [
          ts.createKeywordTypeNode(ts.SyntaxKind.AnyKeyword)
        ]),
        undefined
      )
    }

    function createState(meta: OutputMeta): ts.PropertyDeclaration {
      return ts.createProperty(
        [ts.createDecorator(ts.createCall(ts.createIdentifier(Decorators.State), undefined, []))],
        undefined,
        meta.propDeclarationName,
        undefined,
        meta.typeIdentifier,
        undefined
      )
    }
    function createWatcher(meta: OutputMeta): ts.MethodDeclaration {
      const newValue = 'newValue'
      return ts.createMethod(
        [
          ts.createDecorator(
            ts.createCall(ts.createIdentifier(Decorators.Watch), undefined, [
              ts.createStringLiteral(meta.watchedPropName)
            ])
          )
        ],
        undefined,
        undefined,
        `${meta.emitMethodName}Watcher`,
        undefined,
        undefined,
        [ts.createParameter(undefined, undefined, undefined, newValue, undefined, undefined, undefined)], // parameters
        undefined,
        ts.createBlock(
          [
            ts.createIf(
              ts.createIdentifier(newValue),
              ts.createBlock([
                ts.createStatement(
                  ts.createCall(
                    ts.createPropertyAccess(ts.createPropertyAccess(ts.createThis(), meta.emitMethodName), 'emit'),
                    undefined,
                    [
                      ts.createObjectLiteral([
                        ts.createPropertyAssignment(meta.propDeclarationName, ts.createIdentifier(newValue))
                      ])
                    ]
                  )
                )
              ])
            )
          ],
          true
        )
      )
    }
  }
}

export function outputEventName(prefix: string, suffix?: string): string {
  const _suffix = suffix || 'Saved'
  return `${prefix}${_suffix.charAt(0).toUpperCase() + _suffix.slice(1)}`
}

type OutputMeta = {
  emitMethodName: string
  propDeclarationName: string
  typeIdentifier?: ts.TypeNode
  watchedPropName: string
}
