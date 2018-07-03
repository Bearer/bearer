import * as ts from 'typescript'
import * as d from './declarations'
import { isBearerComponent } from './utils'

const inScope = (components, tag) => {
  for (let component in components) {
    if (components[component].tag === tag) return true
  }
  return false
}

export default (
  metadatas: d.PluginMetadata
): ts.TransformerFactory<ts.SourceFile> => {
  return (transformContext: ts.TransformationContext) => {
    const SCOPE = process.env.TAG_NAMES_SCOPE

    function getNewTagName(tagName: string) {
      if (SCOPE && inScope(metadatas.components, tagName)) {
        return `${SCOPE}-${tagName}`
      }
      return tagName
    }

    function visitJsxElement(node: ts.JsxElement) {
      const { openingElement, closingElement, children } = node
      const oTagName = openingElement.tagName
      openingElement.tagName = ts.createLiteral(
        getNewTagName(oTagName.getText())
      )
      closingElement.tagName = ts.createLiteral(
        getNewTagName(oTagName.getText())
      )
      return ts.updateJsxElement(node, openingElement, children, closingElement)
    }

    function visitJsxSelfClosingElement(node: ts.JsxSelfClosingElement) {
      node.tagName = ts.createLiteral(getNewTagName(node.tagName.getText()))
      return node
    }

    function visitDecorator(node: ts.Decorator) {
      if (isBearerComponent(node)) {
        return ts.visitEachChild(node, visitDecoratorChildren, transformContext)
      }
      return node
    }

    function visitDecoratorChildren(node: ts.Node) {
      switch (node.kind) {
        case ts.SyntaxKind.CallExpression:
          return updateDecoratorCallExpression(node as ts.CallExpression)
      }
      return ts.visitEachChild(node, visitDecoratorChildren, transformContext)
    }

    function updateDecoratorCallExpression(node: ts.CallExpression) {
      // return node
      // we assume @Component is called this way => @Component({...})
      const params: ts.ObjectLiteralExpression = node
        .arguments[0] as ts.ObjectLiteralExpression
      // const property: ts.ObjectLiteralElementLike = params.properties.filter(
      //   p => p.name.getText() === 'tag'
      // )[0]

      const newParams = ts.createObjectLiteral([
        // ...params.properties,
        ts.createPropertyAssignment(
          ts.createLiteral('tag'),
          ts.createLiteral('scoped---')
        )
      ])

      // return ts.updateCall(node, node.expression, node.typeArguments, [
      //   newParams
      // ])

      return ts.updateCall(
        node,
        node.expression,
        node.typeArguments,
        node.arguments
      )
    }

    function visitSourceFile(node: ts.Node): ts.VisitResult<ts.Node> {
      switch (node.kind) {
        case ts.SyntaxKind.JsxElement: {
          // add custom logic here to replace only component within the scr/components folder
          return ts.visitEachChild(
            visitJsxElement(node as ts.JsxElement),
            visitSourceFile,
            transformContext
          )
        }

        case ts.SyntaxKind.JsxSelfClosingElement:
          return ts.visitEachChild(
            visitJsxSelfClosingElement(node as ts.JsxSelfClosingElement),
            visitSourceFile,
            transformContext
          )
        case ts.SyntaxKind.Decorator: {
          return ts.visitEachChild(
            visitDecorator(node as ts.Decorator),
            visitSourceFile,
            transformContext
          )
        }
        default:
      }
      return ts.visitEachChild(node, visitSourceFile, transformContext)
    }

    return (node: ts.SourceFile) => {
      if (node.isDeclarationFile) {
        return node
      }
      return ts.visitEachChild(node, visitSourceFile, transformContext)
    }
  }
}
