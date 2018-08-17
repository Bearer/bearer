/*
 * From this:
 * class AComponent {
 *  @Intent('aName') fetcher: BearerFetch
 * }
 * 
 * to this:
 * 
 * class AComponent {
 *  private fetcher: BearerFetch
 * 
 *  constructor() {
 *    Intent('aName')(this, fetcher)
 *  }
 * }
 *
 * why:
 *  By doing this we allow our Decorator to have access to the component instance instead of the prototype
 */
import * as ts from 'typescript'

import { Decorators } from '../constants'
import { hasDecoratorNamed } from '../helpers/decorator-helpers'
import { TransformerOptions } from '../types'

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

export default function ComponentTransformer({  }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
  return transformContext => {
    // create constructor if it does not exist
    // append Intent call within constructor
    // remove @Intent decorator from the sourcefile

    return tsSourceFile => {
      const registeredIntents: Array<ts.PropertyDeclaration> = []

      const withDecoratorReplaced = visitRemoveIntentDecorators(tsSourceFile as ts.Node, registeredIntents)

      const withConstructor = visitEnsureConstructor(withDecoratorReplaced as ts.Node) as ts.SourceFile

      return visitConstructor(withConstructor as ts.Node, registeredIntents) as ts.SourceFile
    }

    // Remove decorators and replace them with a property access
    function visitRemoveIntentDecorators(
      node: ts.Node,
      registeredIntents: Array<ts.PropertyDeclaration>
    ): ts.VisitResult<ts.Node> {
      if (ts.isPropertyDeclaration(node)) {
        return replaceIfIntentDecorated(node, registeredIntents)
      }
      return ts.visitEachChild(node, node => visitRemoveIntentDecorators(node, registeredIntents), transformContext)
    }

    function replaceIfIntentDecorated(node: ts.PropertyDeclaration, registeredIntents: Array<ts.PropertyDeclaration>) {
      if (hasDecoratorNamed(node, Decorators.Intent)) {
        registeredIntents.push(node)
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

    // Call Intent function

    function visitConstructor(
      node: ts.Node,
      registeredIntents: Array<ts.PropertyDeclaration>
    ): ts.VisitResult<ts.Node> {
      if (ts.isConstructorDeclaration(node)) {
        return addIntentCallToConstructor(node as ts.ConstructorDeclaration, registeredIntents)
      }
      return ts.visitEachChild(node, node => visitConstructor(node, registeredIntents), transformContext)
    }

    function addIntentCallToConstructor(
      node: ts.ConstructorDeclaration,
      registeredIntents: Array<ts.PropertyDeclaration>
    ): ts.Node {
      const intentCalls: Array<ts.Statement> = registeredIntents.map((intent: ts.PropertyDeclaration) => {
        const call: ts.CallExpression = intent.decorators[0].getChildAt(1) as ts.CallExpression
        return ts.createStatement(
          ts.createCall(
            ts.createCall(ts.createIdentifier(Decorators.Intent) as ts.Expression, undefined, call.arguments),
            undefined,
            [ts.createThis(), ts.createLiteral(intent.name.getText())]
          )
        )
      })
      return ts.updateConstructor(
        node,
        node.decorators,
        node.modifiers,
        node.parameters,
        ts.updateBlock(node.body, [...node.body.statements, ...intentCalls])
      )
    }
  }
}
