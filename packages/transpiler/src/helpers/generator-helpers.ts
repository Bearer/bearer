import * as ts from 'typescript'

import { createOrUpdateComponentDidLoad } from '../../src/transformers/bearer'

import { Decorators } from './../constants'
import { CreateFetcherMeta, TAddAutoLoad, TCreateLoadDataCall, TCreateLoadResourceMethod } from './../types'
import { idName } from './name-helpers'

export function createFetcher(meta: CreateFetcherMeta) {
  return ts.createProperty(
    [
      ts.createDecorator(
        ts.createCall(ts.createIdentifier(Decorators.Intent), undefined, [ts.createLiteral(meta.intentName)])
      )
    ],
    undefined,
    meta.intentMethodName,
    undefined,
    undefined,
    undefined
  )
}

function propertyReferenceIdNames(meta: TCreateLoadResourceMethod, metaCollection: TCreateLoadResourceMethod[]) {
  return meta.intentArguments.map(name => {
    const metaInfo = metaCollection.find(meta => meta.propDeclarationName === name)
    if (metaInfo) {
      return ts.createPropertyAssignment(
        idName(name),
        ts.createPropertyAccess(ts.createThis(), metaInfo.propertyReferenceIdName)
      )
    }
    return ts.createPropertyAssignment(name, ts.createPropertyAccess(ts.createThis(), name))
  })
}

export function createLoadResourceMethod(meta: TCreateLoadResourceMethod, metaCollection: TCreateLoadResourceMethod[]) {
  const intentCall = ts.createCall(ts.createPropertyAccess(ts.createThis(), meta.intentMethodName), undefined, [
    ts.createObjectLiteral([
      ts.createPropertyAssignment(
        meta.intentReferenceIdKeyName,
        ts.createPropertyAccess(ts.createThis(), meta.propertyReferenceIdName)
      ),
      ...propertyReferenceIdNames(meta, metaCollection)
    ])
  ])
  const udapteState = ts.createArrowFunction(
    undefined,
    undefined,
    [
      ts.createParameter(
        undefined,
        undefined,
        undefined,
        ts.createObjectBindingPattern([ts.createBindingElement(undefined, undefined, 'data')]),
        undefined,
        ts.createTypeLiteralNode([
          ts.createPropertySignature(undefined, 'data', undefined, meta.typeIdentifier, undefined)
        ]),
        undefined
      )
    ],
    undefined,
    undefined,
    ts.createBlock(
      [
        ts.createIf(
          ts.createIdentifier('data'),
          ts.createBlock([
            ts.createStatement(
              ts.createBinary(
                ts.createPropertyAccess(ts.createThis(), meta.propDeclarationName),
                ts.SyntaxKind.EqualsToken,
                ts.createIdentifier('data')
              )
            )
          ])
        )
      ],
      true
    )
  )
  const promiseHandler = ts.createCall(ts.createPropertyAccess(intentCall, 'then'), undefined, [udapteState])
  return ts.createProperty(
    undefined,
    undefined,
    meta.loadMethodName,
    undefined,
    undefined,
    ts.createArrowFunction(
      undefined,
      undefined,
      undefined,
      undefined,
      undefined,
      ts.createBlock([ts.createStatement(promiseHandler)], true)
    )
  )
}

export function createLoadDataCall(meta: TCreateLoadDataCall) {
  return ts.createStatement(
    ts.createCall(ts.createPropertyAccess(ts.createThis(), meta.loadMethodName), undefined, undefined)
  )
}

export function addAutoLoad(tsClass: ts.ClassDeclaration, meta: TAddAutoLoad): ts.ClassDeclaration {
  if (meta.autoLoad) {
    return createOrUpdateComponentDidLoad(tsClass, block =>
      ts.updateBlock(block, [
        ...block.statements,
        ts.createIf(
          ts.createPropertyAccess(ts.createThis(), meta.propertyReferenceIdName),
          ts.createBlock([createLoadDataCall(meta)])
        )
      ])
    )
  }
  return tsClass
}
