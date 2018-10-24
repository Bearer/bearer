import * as ts from 'typescript'

import Metadata from '../metadata'

type TransformerOptions = {
  verbose?: true
  metadata: Metadata
}

export default function ComponentTagNameScoping({
  metadata
}: TransformerOptions): ts.TransformerFactory<ts.SourceFile> {
  function getNewTagName(tagName: string) {
    const component = metadata.components.find(aComponent => aComponent.initialTagName === tagName)
    return component ? component.finalTagName : tagName
  }

  function visitJsxElement(node: ts.JsxElement) {
    const { openingElement, closingElement, children } = node
    const oTagName = openingElement.tagName.getText()
    const finalTagName = ts.createIdentifier(getNewTagName(oTagName))
    return ts.updateJsxElement(
      node,
      ts.updateJsxOpeningElement(openingElement, finalTagName, openingElement.typeArguments, openingElement.attributes),
      children,
      ts.updateJsxClosingElement(closingElement, finalTagName)
    )
  }

  function visitJsxSelfClosingElement(node: ts.JsxSelfClosingElement) {
    return ts.updateJsxSelfClosingElement(
      node,
      ts.createIdentifier(getNewTagName(node.tagName.getText())),
      node.typeArguments,
      node.attributes
    )
  }

  return _transformContext => {
    function visit(tsNode: ts.Node): ts.VisitResult<ts.Node> {
      switch (tsNode.kind) {
        case ts.SyntaxKind.JsxElement:
          return ts.visitEachChild(visitJsxElement(tsNode as ts.JsxElement), visit, _transformContext)
        case ts.SyntaxKind.JsxSelfClosingElement:
          return ts.visitEachChild(
            visitJsxSelfClosingElement(tsNode as ts.JsxSelfClosingElement),
            visit,
            _transformContext
          )
        default:
      }
      return ts.visitEachChild(tsNode, visit, _transformContext)
    }

    return tsSourceFile => {
      if (tsSourceFile.isDeclarationFile) {
        return tsSourceFile
      }
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
