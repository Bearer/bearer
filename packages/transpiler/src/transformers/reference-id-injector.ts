import * as ts from 'typescript'

import { Decorators, Properties, Types } from '../constants'
import { decoratorNamed, hasPropDecoratedWithName, propDecoratedWithName } from '../helpers/decorator-helpers'
import { TransformerOptions } from '../types'

import { elementDecorator, ensureImportsFromCore, propDecorator } from './bearer'

export default function bearerReferenceIdInjector({  }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return transformContext => {
    return tsSourceFile => {
      if (!hasRetrieveOrSaveStateIntent(tsSourceFile)) {
        return tsSourceFile
      }

      function visit(node: ts.Node): ts.VisitResult<ts.Node> {
        if (ts.isClassDeclaration(node)) {
          return injectHTMLElementPropery(injectBearerReferenceIdProp(node))
        }
        return ts.visitEachChild(node, visit, transformContext)
      }
      return ts.visitEachChild(
        ensureImportsFromCore(tsSourceFile, [Decorators.Prop, Decorators.Element]),
        visit,
        transformContext
      )
    }
  }
}

function injectHTMLElementPropery(tsClass: ts.ClassDeclaration): ts.ClassDeclaration {
  if (hasPropDecoratedWithName(tsClass, Decorators.Element)) {
    const existingProp = propDecoratedWithName(tsClass, Decorators.Element)[0]
    const propertyName = (existingProp.name as ts.Identifier).escapedText as string

    if (propertyName !== Properties.Element) {
      return ts.updateClassDeclaration(
        tsClass,
        tsClass.decorators,
        tsClass.modifiers,
        tsClass.name,
        tsClass.typeParameters,
        tsClass.heritageClauses,
        [
          ...tsClass.members,
          ts.createGetAccessor(
            undefined, // decorators
            undefined, // modifiers
            Properties.Element, // name
            undefined, // questionExclamationToken
            ts.createTypeReferenceNode(Types.HTMLElement, undefined), // Type
            ts.createBlock([
              ts.createReturn(
                ts.createPropertyAccess(ts.createThis(), propertyName) // initializer
              )
            ])
          )
        ]
      )
    }
    return tsClass
  }

  return ts.updateClassDeclaration(
    tsClass,
    tsClass.decorators,
    tsClass.modifiers,
    tsClass.name,
    tsClass.typeParameters,
    tsClass.heritageClauses,
    [
      ...tsClass.members,
      ts.createProperty(
        [elementDecorator()],
        undefined,
        Properties.Element,
        undefined,
        ts.createTypeReferenceNode(Types.HTMLElement, undefined),
        undefined
      )
    ]
  )
}

function injectBearerReferenceIdProp(tsClass: ts.ClassDeclaration): ts.ClassDeclaration {
  return ts.updateClassDeclaration(
    tsClass,
    tsClass.decorators,
    tsClass.modifiers,
    tsClass.name,
    tsClass.typeParameters,
    tsClass.heritageClauses,
    [
      ...tsClass.members,
      ts.createProperty(
        [propDecorator()],
        undefined,
        Properties.ReferenceId,
        undefined,
        ts.createKeywordTypeNode(ts.SyntaxKind.StringKeyword),
        undefined
      )
    ]
  )
}

function hasRetrieveOrSaveStateIntent(tsSourceFile: ts.SourceFile): boolean {
  return ts.forEachChild(
    tsSourceFile,
    (node): boolean => {
      if (ts.isClassDeclaration(node)) {
        return hasIntentWithStateTypeDecorator(node)
      }
      return false
    }
  )
}

function hasIntentWithStateTypeDecorator(classNode: ts.ClassDeclaration): boolean {
  const properties = propDecoratedWithName(classNode, Decorators.Intent)

  return properties.reduce((has, prop) => {
    const decorator = prop.decorators.find(deco => decoratorNamed(deco, Decorators.Intent))
    if (!decorator || !ts.isCallExpression(decorator.expression)) {
      return has
    }
    const intentType = (decorator.expression as ts.CallExpression).arguments[1] as ts.PropertyAccessExpression
    if (!intentType) {
      return false
    }
    const hasWisthState =
      (intentType.expression as ts.Identifier).escapedText.toString() === Types.IntentType &&
      (intentType.name as ts.Identifier).escapedText.toString() === Types.SaveState
    return has || hasWisthState
  }, false)
}
