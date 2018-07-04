import * as ts from 'typescript'
import * as d from './declarations'
import { isBearerComponent } from './utils'

function filterParams(
  params: ts.ObjectLiteralExpression
): ts.ObjectLiteralExpression {
  // return ts.createObjectLiteral(
  // console.log(
  //   '[BEARER]',
  //   'params.properties',
  //   // params.properties,
  //   params.properties.map(p => p.name.getText())
  // )

  return params
  // return ts.updateObjectLiteral(
  //   params,
  //   params.properties.filter(p => {
  //     if (!p.name)
  //       console.log("Cédric: ", p)
  //     return p.name && p.name.getText() !== 'bearer'
  //   })
  // )
}

export default function DecoratorTransformer({
  verbose
}: d.PluginOptions = {}): ts.TransformerFactory<ts.SourceFile> {
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

    const propDecorator = ts.createDecorator(
      ts.createCall(
        ts.createIdentifier('Prop') as ts.Expression,
        undefined,
        undefined
      )
    )

    function visitClass(classNode: ts.ClassDeclaration): ts.ClassDeclaration {
      console.log(classNode.name)
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
            [propDecorator],
            undefined,
            'BEARER_ID',
            undefined,
            undefined,
            undefined
          ),
          ts.createProperty(
            undefined,
            [],
            'SCENARIO_ID',
            undefined,
            undefined,
            ts.createStringLiteral('BEARER_SCENARIO_ID')
          )
        ]
      )
    }

    function visit(node: ts.Node): ts.VisitResult<ts.Node> {
      switch (node.kind) {
        case ts.SyntaxKind.Decorator: {
          const decoNode = node as ts.Decorator
          if (isBearerComponent(node as ts.Decorator)) {
            return visitDecorator(node as ts.Decorator)
          } else {
            return ts.visitEachChild(node, visit, transformContext)
          }
        }
        case ts.SyntaxKind.ClassDeclaration: {
          return ts.visitEachChild(
            visitClass(node as ts.ClassDeclaration),
            visit,
            transformContext
          )
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
