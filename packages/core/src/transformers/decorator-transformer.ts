import * as ts from 'typescript'

function filterParams(
  params: ts.ObjectLiteralExpression
): ts.ObjectLiteralExpression {
  return ts.createObjectLiteral(
    params.properties.filter(p => {
      return p.name.getText() !== 'bearer'
    })
  )
}

type TransformerOptions = {
  verbose?: true
}
export default function DecoratorTransdormer({
  verbose
}: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  function log(...args) {
    if (verbose) {
      console.log.apply(this, args)
    }
  }

  return transformContext => {
    log('[BEARER]', 'Using Decorator Transformer')
    function visitDecoratorCall(node: ts.CallExpression) {
      return ts.updateCall(node, node.expression, node.typeArguments, [
        filterParams(node.arguments[0] as ts.ObjectLiteralExpression)
      ])
    }

    function visitDecorator(node: ts.Node) {
      switch (node.kind) {
        case ts.SyntaxKind.CallExpression:
          return visitDecoratorCall(node as ts.CallExpression)
      }
      return ts.visitEachChild(node, visitDecorator, transformContext)
    }

    function visitClass(classNode: ts.ClassDeclaration): ts.ClassDeclaration {
      return ts.updateClassDeclaration(
        classNode,
        classNode.decorators,
        classNode.modifiers,
        classNode.name,
        classNode.typeParameters,
        classNode.heritageClauses,
        [
          ...classNode.members,
          ts.createGetAccessor(
            undefined,
            [],
            'BEARER_ID',
            undefined,
            undefined,
            ts.createBlock([
              ts.createReturn(ts.createIdentifier('this.SCENARIO_ID'))
            ])
          ),
          ts.createGetAccessor(
            undefined,
            [],
            'SCENARIO_ID',
            undefined,
            undefined,
            ts.createBlock([
              ts.createReturn(ts.createIdentifier('"BEARER_SCENARIO_ID"'))
            ])
          )
        ]
      )
    }

    function visit(node: ts.Node): ts.VisitResult<ts.Node> {
      switch (node.kind) {
        case ts.SyntaxKind.ClassDeclaration: {
          return ts.visitEachChild(
            visitClass(node as ts.ClassDeclaration),
            visit,
            transformContext
          )
        }
        case ts.SyntaxKind.Decorator: {
          const decoNode = node as ts.Decorator
          log(
            '[BEARER]',
            'decorator',
            decoNode.expression.getChildCount() &&
              decoNode.getChildAt(1).getText()
          )
          if (
            decoNode.expression.getChildCount() &&
            decoNode
              .getChildAt(0)
              .getText()
              .match(/^Component/)
          ) {
            return visitDecorator(node as ts.Decorator)
          }
        }
      }
      return ts.visitEachChild(node, visit, transformContext)
    }

    return tsSourceFile => {
      log('[BEARER]', 'exploring', tsSourceFile.fileName)
      return visit(tsSourceFile) as ts.SourceFile
    }
  }
}
