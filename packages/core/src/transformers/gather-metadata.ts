import * as ts from 'typescript'
import * as d from './declarations'
import { isBearerComponent, getComponentDecoratorTagName } from './utils'

const find = (ary, fun) => {
  for (let el in ary) {
    let obj = ary[el]
    if (fun(obj)) return obj
  }
}
function gatherMetadata(
  metadatas: d.PluginMetadata
): ts.TransformerFactory<ts.SourceFile> {
  return transformContext => {
    function visitClassDeclaration(nodeClass: ts.ClassDeclaration): void {
      //add component inforamtion to metadata
      if (nodeClass.decorators) {
        const ComponentDecorator = find(nodeClass.decorators, isBearerComponent)
        if (ComponentDecorator) {
          metadatas.components.push({
            tag: getComponentDecoratorTagName(ComponentDecorator)
          })
        }
      }
    }

    function visitNode(node: ts.Node): ts.VisitResult<ts.Node> {
      switch (node.kind) {
        case ts.SyntaxKind.ClassDeclaration: {
          visitClassDeclaration(node as ts.ClassDeclaration)
          return node
        }
      }

      return ts.visitEachChild(node, visitNode, transformContext)
    }

    return (tsSourceFile: ts.SourceFile) => {
      if (tsSourceFile.isDeclarationFile) {
        return tsSourceFile
      }
      return ts.visitEachChild(tsSourceFile, visitNode, transformContext)
    }
  }
}

export default gatherMetadata
