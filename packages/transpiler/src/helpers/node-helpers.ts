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

export function hasBooleanProperty(
  node: ts.ObjectLiteralExpression,
  prop: string,
  kind: ts.SyntaxKind.TrueKeyword | ts.SyntaxKind.FalseKeyword
) {
  const argument = node.properties[0]
  let identifier: any = { escapedText: 'undefined' }
  let initializer: any = { kind: ts.SyntaxKind.UndefinedKeyword }
  if (argument.kind === ts.SyntaxKind.PropertyAssignment) {
    identifier = argument.name as ts.Identifier
    initializer = argument.initializer
  }
  return (identifier.escapedText || identifier.text) === prop && initializer.kind === kind
}
