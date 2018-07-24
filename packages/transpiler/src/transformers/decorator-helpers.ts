import * as ts from 'typescript'

export function name(node: ts.Decorator): string {
  let name: any = ''
  if (node.getChildAt(1).getChildCount() === 0) {
    name = node.getChildAt(1).getText()
  } else {
    name = node
      .getChildAt(1)
      .getChildAt(0)
      .getText()
  }
  return name
}

export function hasName(node: ts.Decorator, tsDecoratorName: string) {
  return name(node) === tsDecoratorName
}

export function classDecoratedWithName(node: ts.ClassDeclaration, decoratorName: string): boolean {
  let hasComponentDecorator = false
  ts.forEachChild(node, n => {
    if (n.kind === ts.SyntaxKind.Decorator) {
      hasComponentDecorator = hasComponentDecorator || name(n as ts.Decorator) === decoratorName
    }
  })
  return hasComponentDecorator
}

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

export function hasDecoratorNamed(propDeclaration: ts.PropertyDeclaration, name: string): boolean {
  let has = false
  ts.forEachChild(propDeclaration, anode => {
    if (ts.isDecorator(anode)) {
      ts.forEachChild(anode, node => {
        if (ts.isCallExpression(node) && node.expression['escapedText'] === name) {
          has = true
        }
      })
    }
  })
  return has
}

export function using(node: ts.Node, decoratorName: string): boolean {
  let usedInCode = false
  function visit(node: ts.Node) {
    if (ts.isDecorator(node)) usedInCode = usedInCode || hasName(node as ts.Decorator, decoratorName)
    ts.forEachChild(node, visit)
  }

  ts.forEachChild(node, visit)

  return usedInCode
}

export default {
  name,
  classDecoratedWithName,
  hasName,
  using
}
