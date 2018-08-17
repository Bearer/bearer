import * as ts from 'typescript'

import { Decorators } from '../constants'

import { decoratorNamed } from './decorator-helpers'

export function isWatcherOn(tsMethod: ts.MethodDeclaration, watchedProp: string): boolean {
  if (!tsMethod.decorators) {
    return false
  }
  const decorator = tsMethod.decorators.find(d => decoratorNamed(d, Decorators.Watch))
  if (!decorator) {
    return false
  }
  if (!ts.isCallExpression(decorator.expression)) {
    return false
  }
  const firstArgument = decorator.expression.arguments[0]
  return ts.isStringLiteral(firstArgument) && firstArgument.text === watchedProp
}
