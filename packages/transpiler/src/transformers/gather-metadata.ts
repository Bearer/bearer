import * as Case from 'case'
import * as ts from 'typescript'

import { Component, Decorators } from '../constants'
import { getDecoratorNamed, getExpressionFromDecorator, hasDecoratorNamed } from '../helpers/decorator-helpers'
import { TransformerOptions } from '../types'

type TInput = {
  name: string
  type: 'number' | 'string'
  default: string
}

type TOutput = {
  name: string
  payloadFormat: {
    [key: string]: string
  }
}

function propAsInput(tsProp: ts.PropertyDeclaration): TInput {
  return {
    name: (tsProp.name as ts.Identifier).escapedText.toString(),
    type: tsProp.type.kind === ts.SyntaxKind.NumberKeyword ? 'number' : 'string',
    default: (tsProp.initializer as ts.Expression).getText()
  }
}

function eventAsOutput(tsProp: ts.PropertyDeclaration): TOutput {
  return {
    name: (tsProp.name as ts.Identifier).escapedText.toString(),
    payloadFormat: {}
  }
}

function collectInputs(tsClass: ts.ClassDeclaration): Array<TInput> {
  return tsClass.members
    .filter(member => ts.isPropertyDeclaration(member) && hasDecoratorNamed(member, Decorators.Prop))
    .map(propAsInput)
    .filter(prop => prop.name !== Component.bearerContext)
}

function collectOutputs(tsClass: ts.ClassDeclaration): Array<TOutput> {
  return tsClass.members
    .filter(member => ts.isPropertyDeclaration(member) && hasDecoratorNamed(member, Decorators.Event))
    .map(eventAsOutput)
}

export default function GatherMetadata({ metadata }: TransformerOptions): ts.TransformerFactory<ts.SourceFile> {
  function getTagNames(tagName: string): { initialTagName: string; finalTagName: string } {
    const finalTag =
      metadata.prefix && metadata.suffix
        ? [Case.kebab(metadata.prefix), tagName, Case.kebab(metadata.suffix)].join('-')
        : tagName
    return {
      initialTagName: tagName,
      finalTagName: finalTag
    }
  }

  function visit(node: ts.Node): ts.Node {
    // Found Component
    if (ts.isClassDeclaration(node)) {
      if (hasDecoratorNamed(node, Decorators.Component)) {
        const component = getDecoratorNamed(node, Decorators.Component)
        const tag = getExpressionFromDecorator<ts.StringLiteral>(component, 'tag')
        metadata.registerComponent({
          classname: node.name.text,
          isRoot: false,
          ...getTagNames(tag.text)
        })
      }
      // Found RootComponent
      else if (hasDecoratorNamed(node, Decorators.RootComponent)) {
        const component = getDecoratorNamed(node, Decorators.RootComponent)
        const nameExpression = getExpressionFromDecorator<ts.StringLiteral>(component, 'role')
        const name = nameExpression ? nameExpression.text : ''
        const groupExpression = getExpressionFromDecorator<ts.StringLiteral>(component, 'group')
        const group = groupExpression ? groupExpression.text : ''
        const tag = [Case.kebab(group), name].join('-')
        metadata.registerComponent({
          classname: node.name.text,
          isRoot: true,
          ...getTagNames(tag),
          group,
          inputs: collectInputs(node),
          outputs: collectOutputs(node)
        })
      }
    }
    return node
  }

  return _transformContext => {
    return tsSourceFile => {
      return ts.visitEachChild(tsSourceFile, visit, _transformContext)
    }
  }
}
