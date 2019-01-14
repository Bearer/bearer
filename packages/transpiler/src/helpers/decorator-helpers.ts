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

export function getDecoratorProperties(tsDecorator: ts.Decorator, index = 0): ts.ObjectLiteralExpression {
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
  const found: ts.PropertyAssignment = tsLiteral.properties.find(property => {
    if (ts.isStringLiteral(property.name)) {
      return property.name.text === key
    } else {
      return (property.name as ts.Identifier).escapedText.toString() === key
    }
  }) as ts.PropertyAssignment
  if (found) {
    return found.initializer as T
  }
  return null
}

export function getExpressionFromDecorator<T extends ts.Expression>(
  tsDecorator: ts.Decorator,
  key: string,
  index = 0
): T {
  return getExpressionFromLiteralObject(getDecoratorProperties(tsDecorator, index), key)
}

export function getDecoratorNamed(
  tsDecoratedNode: { decorators?: ts.NodeArray<ts.Decorator> },
  name: string
): ts.Decorator | undefined {
  if (tsDecoratedNode.decorators) {
    return tsDecoratedNode.decorators.find(decorator => decoratorNamed(decorator, name))
  }
}

// tslint:disable-next-line:prefer-array-literal
export function extractStringOptions<T>(objectLitteral: ts.ObjectLiteralExpression, keys: Array<keyof T>): Partial<T> {
  const values: Partial<T> = {}
  keys.map(value => {
    const optionValue = getExpressionFromLiteralObject<ts.StringLiteral>(objectLitteral, value as string)
    if (optionValue) {
      values[value as string] = optionValue && optionValue.text
    }
  })
  return values
}

export function extractArrayOptions<T>(
  objectLitteral: ts.ObjectLiteralExpression,
  keys: (keyof T)[]
): { [P in keyof T]?: string[] } {
  const values = {}
  keys.map(value => {
    const optionValue = getExpressionFromLiteralObject<ts.StringLiteral>(objectLitteral, value as string)
    if (optionValue && ts.isArrayLiteralExpression(optionValue)) {
      values[value as string] = optionValue.elements.map(e => {
        if (ts.isStringLiteral(e)) {
          return e.text
        }
        return []
      })
    }
  })
  return values
}

export function extractBooleanOptions<T>(objectLitteral: ts.ObjectLiteralExpression, keys: (keyof T)[]): Partial<T> {
  const values: Partial<T> = {}
  keys.map(value => {
    const optionValue = getExpressionFromLiteralObject<ts.BooleanLiteral>(objectLitteral, value as string)
    if (optionValue) {
      values[value as string] = optionValue && optionValue.kind === ts.SyntaxKind.TrueKeyword
    }
  })

  return values
}
