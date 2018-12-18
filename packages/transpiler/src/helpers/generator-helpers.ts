import * as ts from 'typescript'

import { createOrUpdateComponentDidLoad } from '../../src/transformers/bearer'

import { Decorators } from './../constants'
import { CreateFetcherMeta, TAddAutoLoad, TCreateLoadDataCall, TCreateLoadResourceMethod } from './../types'

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

export function createLoadResourceMethod(meta: TCreateLoadResourceMethod) {
  const intentCall = ts.createCall(ts.createPropertyAccess(ts.createThis(), meta.intentMethodName), undefined, [
    ts.createObjectLiteral([
      ts.createPropertyAssignment(
        meta.intentReferenceIdKeyName,
        ts.createPropertyAccess(ts.createThis(), meta.propertyReferenceIdName)
      )
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
        ts.createStatement(
          ts.createBinary(
            ts.createPropertyAccess(ts.createThis(), meta.propDeclarationName),
            ts.SyntaxKind.EqualsToken,
            ts.createIdentifier('data')
          )
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
