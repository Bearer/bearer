import * as ts from 'typescript'

export function hasPropDecoratedWithName(classNode: ts.ClassDeclaration, decoratorName: string): boolean {
  return Boolean(propDecoratedWithName(classNode, decoratorName).length)
}

export function propDecoratedWithName(
  node: ts.ClassDeclaration,
  decoratorName: string
): Array<ts.PropertyDeclaration | null> {
  const props: Array<ts.PropertyDeclaration | null> = []
  ts.forEachChild(node, node => {
    if (ts.isPropertyDeclaration(node) && hasDecoratorNamed(node as ts.PropertyDeclaration, decoratorName)) {
      props.push(node)
    }
  })
  return props
}

export function hasDecoratorNamed(
  decoratedNode: ts.PropertyDeclaration | ts.ClassDeclaration | ts.MethodDeclaration,
  name: string
): boolean {
  return ts.forEachChild(decoratedNode, anode => {
    return ts.isDecorator(anode) && decoratorNamed(anode, name)
  })
}

export function decoratorNamed(tsDecorator: ts.Decorator, name: string): boolean {
  return ts.forEachChild(tsDecorator, node => {
    return ts.isCallExpression(node) && (node.expression as ts.Identifier).escapedText === name
  })
}

export function getDecoratorProperties(tsDecorator: ts.Decorator, index: number = 0): ts.ObjectLiteralExpression {
  if (!ts.isCallExpression(tsDecorator.expression)) {
    return ts.createObjectLiteral([])
  }
  const expression = tsDecorator.expression
  const argument = expression.arguments[index]

  if (!ts.isObjectLiteralExpression(argument)) {
    return ts.createObjectLiteral([])
  }
  return ts.createObjectLiteral(argument.properties)
}

export function getExpressionFromLiteralObject<T extends ts.Expression>(
  tsLiteral: ts.ObjectLiteralExpression,
  key: string
): T {
  const found: ts.PropertyAssignment = tsLiteral.properties.find(
    property => property.name.getText().trim() === key
  ) as ts.PropertyAssignment
  if (found) {
    return found.initializer as T
  }
  return null
}

export function getExpressionFromDecorator<T extends ts.Expression>(
  tsDecorator: ts.Decorator,
  key: string,
  index: number = 0
): T {
  return getExpressionFromLiteralObject(getDecoratorProperties(tsDecorator, index), key)
}

export function getDecoratorNamed(tsClassNode: ts.ClassDeclaration, name: string): ts.Decorator | undefined {
  return tsClassNode.decorators.find(decorator => decoratorNamed(decorator, name))
}
