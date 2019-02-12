/*
 * scope i18n component and helpers
 */
import * as ts from 'typescript'

// At the same time it ensures we will have the accessor present ;-)
import { shouldProcessFile as hasAccessor, retrieveScenarioId } from './scenario-id-accessor-injector'

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

export default function i18nModifier(_options: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    function normalizeImports(collection: AliasesCollection) {
      function transformNamedBindings(tsImportSpecifier: ts.ImportSpecifier): ts.ImportSpecifier {
        const importName = (tsImportSpecifier.propertyName || tsImportSpecifier.name) as ts.Identifier
        switch (importName.escapedText) {
          case 't': {
            const helperName = 'scopedT'
            collection.push({
              helperName,
              name: tsImportSpecifier.name.escapedText.toString()
            })
            return ts.createImportSpecifier(ts.createIdentifier(helperName), tsImportSpecifier.name)
          }
          case 'p': {
            const helperName = 'scopedP'
            collection.push({
              helperName,
              name: tsImportSpecifier.name.escapedText.toString()
            })
            return ts.createImportSpecifier(ts.createIdentifier(helperName), tsImportSpecifier.name)
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

      return (tsNode: ts.Node): ts.VisitResult<ts.Node> => {
        switch (tsNode.kind) {
          // alias imports
          case ts.SyntaxKind.ImportDeclaration: {
            return ts.visitEachChild(
              aliasImports(tsNode as ts.ImportDeclaration),
              normalizeImports(collection),
              _transformContext
            )
          }
          // if bearer-i18n tag add scope
          case ts.SyntaxKind.JsxSelfClosingElement:
            return ts.visitEachChild(
              visitJsxSelfClosingElement(tsNode as ts.JsxSelfClosingElement),
              normalizeImports(collection),
              _transformContext
            )
        }
        return ts.visitEachChild(tsNode, normalizeImports(collection), _transformContext)
      }
    }

    return tsSourceFile => {
      if (!hasAccessor(tsSourceFile)) {
        return tsSourceFile
      }
      const aliasCollection: AliasesCollection = []
      const normalizedImports = ts.visitEachChild(tsSourceFile, normalizeImports(aliasCollection), _transformContext)
      if (!aliasCollection.length) {
        return normalizedImports
      }

      return injectScopedInstances(normalizedImports, aliasCollection)
    }
  }
}

function injectScopedInstances(tsSourceFile: ts.SourceFile, collection: AliasesCollection): ts.SourceFile {
  return ts.updateSourceFileNode(
    tsSourceFile,
    [
      ...tsSourceFile.statements,
      ...collection.map(({ name, helperName }) =>
        ts.createVariableStatement(
          undefined, // no using this [ts.createToken(ts.SyntaxKind.ConstKeyword)], otherwise it produces : const var t
          [
            ts.createVariableDeclaration(
              name,
              undefined,
              ts.createCall(ts.createIdentifier(helperName), undefined, [ts.createStringLiteral(retrieveScenarioId())])
            )
          ]
        )
      )
    ],
    tsSourceFile.isDeclarationFile,
    tsSourceFile.referencedFiles,
    tsSourceFile.referencedFiles,
    tsSourceFile.hasNoDefaultLib,
    tsSourceFile.libReferenceDirectives
  )
}

type AliasesCollection = {
  name: string
  helperName: 'scopedT' | 'scopedP'
}[]
