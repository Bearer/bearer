import * as ts from 'typescript'

export function hasMethodNamed(classNode: ts.ClassDeclaration, methodName: string): boolean {
  return ts.forEachChild(classNode, node => isMethodNamed(node, methodName))
}

export function isMethodNamed(tsNode: ts.Node, name: string): boolean {
  return ts.isMethodDeclaration(tsNode) && (tsNode.name as ts.Identifier).escapedText === name
}

export function getNodeName(tsNode: { name: ts.PropertyName }): string {
  return (tsNode.name as ts.Identifier).escapedText.toString()
}
