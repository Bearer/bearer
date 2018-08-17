import * as ts from 'typescript'

import { hasMethodNamed } from './node-helpers'

export function ensureMethodExists(classNode: ts.ClassDeclaration, methodName: string) {
  if (hasMethodNamed(classNode, methodName)) {
    return classNode
  }
  return ts.updateClassDeclaration(
    classNode,
    classNode.decorators,
    classNode.modifiers,
    classNode.name,
    classNode.typeParameters,
    classNode.heritageClauses,
    [
      ...classNode.members,
      ts.createMethod(
        undefined,
        undefined,
        undefined,
        methodName,
        undefined,
        undefined,
        undefined,
        undefined,
        ts.createBlock([])
      )
    ]
  )
}
