/*
 * From this:
 * class AComponent {
 *  @Function('aName') fetcher: BearerFetch
 * }
 *
 * to this:
 *
 * class AComponent {
 *  private fetcher: BearerFetch
 *
 *  constructor() {
 *    Function('aName')(this, fetcher)
 *  }
 * }
 *
 * why:
 *  By doing this we allow our Decorator to have access to the component instance instead of the prototype
 */
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { hasDecoratorNamed } from '../helpers/decorator-helpers'
import { getNodeName } from '../helpers/node-helpers'
import { TransformerOptions } from '../types'
import debug from '../logger'

const logger = debug.extend('replace-function-decorators')
function appendConstructor(node: ts.ClassDeclaration): ts.Node {
  if (classHasConstructor(node)) {
    return node
  }
  return ts.updateClassDeclaration(
    node,
    node.decorators,
    node.modifiers,
    node.name,
    node.typeParameters,
    node.heritageClauses,
    [
      ...node.members,
      ts.createConstructor(
        /* constructors */ undefined,
        /* modifiers */ undefined,
        /* parameters */ undefined,
        ts.createBlock([])
      )
    ]
  )
}

function classHasConstructor(classNode: ts.ClassDeclaration): boolean {
  return ts.forEachChild(classNode, aNode => {
    return ts.isConstructorDeclaration(aNode)
  })
}

export default function componentTransformer({  }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return transformContext => {
    // create constructor if it does not exist
    // append Function call within constructor
    // remove @Function decorator from the sourcefile

    return tsSourceFile => {
      logger('processing %s', tsSourceFile.fileName)
      const registeredFunctions: ts.PropertyDeclaration[] = []

      const withDecoratorReplaced = visitRemoveFunctionDecorators(tsSourceFile as ts.Node, registeredFunctions)

      const withConstructor = visitEnsureConstructor(withDecoratorReplaced as ts.Node) as ts.SourceFile

      return visitConstructor(withConstructor as ts.Node, registeredFunctions) as ts.SourceFile
    }

    // Remove decorators and replace them with a property access
    function visitRemoveFunctionDecorators(
      node: ts.Node,
      registeredFunctions: ts.PropertyDeclaration[]
    ): ts.VisitResult<ts.Node> {
      if (ts.isPropertyDeclaration(node)) {
        return replaceIfFunctionDecorated(node, registeredFunctions)
      }
      return ts.visitEachChild(node, node => visitRemoveFunctionDecorators(node, registeredFunctions), transformContext)
    }

    function replaceIfFunctionDecorated(node: ts.PropertyDeclaration, registeredFunctions: ts.PropertyDeclaration[]) {
      if (hasDecoratorNamed(node, Decorators.Function)) {
        registeredFunctions.push(node)
        return ts.createProperty(
          /* decorators */ undefined,
          /* modifiers */ [ts.createToken(ts.SyntaxKind.PrivateKeyword)],
          node.name,
          node.questionToken,
          node.type,
          node.initializer
        )
      }
      return node
    }
    // Create a constructor if none is present

    function visitEnsureConstructor(node: ts.Node): ts.VisitResult<ts.Node> {
      if (ts.isClassDeclaration(node)) {
        return ts.visitEachChild(
          appendConstructor(node as ts.ClassDeclaration),
          visitEnsureConstructor,
          transformContext
        )
      }
      return ts.visitEachChild(node, visitEnsureConstructor, transformContext)
    }

    // Call Function function

    function visitConstructor(node: ts.Node, registeredFunctions: ts.PropertyDeclaration[]): ts.VisitResult<ts.Node> {
      if (ts.isConstructorDeclaration(node)) {
        return addFunctionCallToConstructor(node as ts.ConstructorDeclaration, registeredFunctions)
      }
      return ts.visitEachChild(node, node => visitConstructor(node, registeredFunctions), transformContext)
    }

    function addFunctionCallToConstructor(
      node: ts.ConstructorDeclaration,
      registeredFunctions: ts.PropertyDeclaration[]
    ): ts.Node {
      const functionCalls: ts.Statement[] = registeredFunctions.map((func: ts.PropertyDeclaration) => {
        const call: ts.CallExpression = func.decorators[0].expression as ts.CallExpression
        return ts.createStatement(
          ts.createCall(
            ts.createCall(ts.createIdentifier(Decorators.BackendFunction) as ts.Expression, undefined, call.arguments),
            undefined,
            [ts.createThis(), ts.createLiteral(getNodeName(func))]
          )
        )
      })
      return ts.updateConstructor(
        node,
        node.decorators,
        node.modifiers,
        node.parameters,
        ts.createBlock([...node.body.statements, ...functionCalls], true)
      )
    }
  }
}
