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

function propDecorator() {
  return ts.createDecorator(
    ts.createCall(
      ts.createIdentifier('Prop') as ts.Expression,
      undefined,
      undefined
    )
  )
}

type TransformerOptions = {
  verbose?: true
}
export default function ComponentTransformer({
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
          ts.createProperty(
            [propDecorator()],
            undefined,
            'BEARER_ID',
            undefined,
            ts.createKeywordTypeNode(ts.SyntaxKind.StringKeyword),
            ts.createStringLiteral('BEARER_SCENARIO_ID')
          ),
          ts.createProperty(
            [propDecorator()],
            undefined,
            'SCENARIO_ID',
            undefined,
            ts.createKeywordTypeNode(ts.SyntaxKind.StringKeyword),
            ts.createStringLiteral('BEARER_SCENARIO_ID')
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
          if (isComponentDecorator(node as ts.Decorator)) {
            return visitDecorator(node as ts.Decorator)
          } else {
            return ts.visitEachChild(node, visit, transformContext)
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

function isComponentDecorator(node: ts.Decorator): boolean {
  return (
    node.expression &&
    node.expression['expression'] &&
    node.expression['expression'].escapedText === 'Component'
  )
}
