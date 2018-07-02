import * as ts from 'typescript'

export function isDecoratorNamed(name: string) {
  return (decorator: ts.Decorator): boolean => {
    if (decorator.expression.getChildCount()) {
      return decorator.expression.getChildren()[0].getText() === name
    }
    return false
  }
}

export const isBearerComponent = isDecoratorNamed('Component')

export function getComponentDecoratorTagName(decorator: ts.Decorator): string {
  const call: ts.CallExpression = decorator.getChildAt(1) as ts.CallExpression
  const params: ts.ObjectLiteralExpression = call
    .arguments[0] as ts.ObjectLiteralExpression
  const property: ts.ObjectLiteralElementLike = params.properties.filter(
    p => p.name.getText() === 'tag'
  )[0]

  // TODO: there is probably a better way to get the string value
  return property
    ? property
        .getLastToken()
        .getText()
        .replace(/['\"]/g, '')
    : 'MISSING_TAG_NAME'
}
