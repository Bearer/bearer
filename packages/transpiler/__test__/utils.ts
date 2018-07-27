import * as ts from 'typescript'
export default (sourceFile: ts.SourceFile) => {
  let names: string[] = []
  traverseNode(sourceFile, names)
  function traverseNode(node: ts.Node, names: string[]) {
    if (node.kind === ts.SyntaxKind.Decorator) names.push(node.expression.expression.escapedText)
    ts.forEachChild(node, node => {
      traverseNode(node, names)
    })
  }
  return names
}
