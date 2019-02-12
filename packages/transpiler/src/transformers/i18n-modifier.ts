/*
 * scope i18n component and helpers
 */
import * as ts from 'typescript'

// At the same time it ensures we will have the accessor present ;-)
import { shouldProcessFile as hasAccessor } from './scenario-id-accessor-injector'

import { TransformerOptions } from '../types'
import { Component, Module } from '../constants'

const I18N_TAG = 'bearer-i18n'

function visitJsxSelfClosingElement(node: ts.JsxSelfClosingElement): ts.JsxSelfClosingElement {
  if ((node.tagName as ts.Identifier).escapedText.toString().toLowerCase() !== I18N_TAG) {
    return node
  }
  return ts.updateJsxSelfClosingElement(
    node,
    node.tagName,
    node.typeArguments,
    ts.createJsxAttributes([
      ...node.attributes.properties,
      ts.createJsxAttribute(
        ts.createIdentifier('scope'),
        ts.createJsxExpression(undefined, ts.createPropertyAccess(ts.createThis(), Component.scenarioIdAccessor))
      )
    ])
  )
}

function transformNamedBindings(tsImportSpecifier: ts.ImportSpecifier): ts.ImportSpecifier {
  const importName = (tsImportSpecifier.propertyName || tsImportSpecifier.name) as ts.Identifier
  switch (importName.escapedText) {
    case 't': {
      return ts.createImportSpecifier(ts.createIdentifier('scopedT'), tsImportSpecifier.name)
    }
    case 'p': {
      return ts.createImportSpecifier(ts.createIdentifier('scopedP'), tsImportSpecifier.name)
    }
  }
  return tsImportSpecifier
}

function aliasImports(tsImport: ts.ImportDeclaration): ts.ImportDeclaration {
  if ((tsImport.moduleSpecifier as ts.StringLiteral).text.toString() !== Module.BEARER_CORE_MODULE) {
    return tsImport
  }
  const bindings = tsImport.importClause.namedBindings as ts.NamedImports

  return ts.updateImportDeclaration(
    tsImport,
    tsImport.decorators,
    tsImport.modifiers,
    ts.createImportClause(
      tsImport.importClause.name,
      ts.createNamedImports(bindings.elements.map(transformNamedBindings))
    ),
    tsImport.moduleSpecifier
  )
}

export default function i18nModifier(_options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      switch (tsNode.kind) {
        // alias imports
        case ts.SyntaxKind.ImportDeclaration: {
          return ts.visitEachChild(aliasImports(tsNode as ts.ImportDeclaration), visit, _transformContext)
        }
        // create
        // if bearer-i18n tag add scope
        case ts.SyntaxKind.JsxSelfClosingElement:
          return ts.visitEachChild(
            visitJsxSelfClosingElement(tsNode as ts.JsxSelfClosingElement),
            visit,
            _transformContext
          )
      }
      return ts.visitEachChild(tsNode, visit, _transformContext)
    }

    return tsSourceFile => {
      if (!hasAccessor(tsSourceFile)) {
        return tsSourceFile
      }

      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
