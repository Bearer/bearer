import * as ts from 'typescript'

import { Decorators } from './../constants'
import { CreateFetcherMeta } from './../types'

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
