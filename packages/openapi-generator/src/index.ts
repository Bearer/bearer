import fs from 'fs'
import Mustache from 'mustache'
import path from 'path'
import ts from 'typescript'

import { convertType } from './convert-type'

export const INTENT_ACTION = 'action'
type TOutput = { requestBody: any; response: any; intentAuthType?: string }

/**
 * find action method in class and return relevant node
 */
function findActionMethod(node: ts.ClassDeclaration): ts.MethodDeclaration | undefined {
  return node.members.find((node: ts.ClassElement) => {
    if (node.name) {
      return node.name.getText() === INTENT_ACTION
    }
    return false
  }) as ts.MethodDeclaration
}

/**
 *  Convert intent types to json schemas
 */

function serializeParameters(sym: ts.Symbol | undefined, node: ts.Node, checker: ts.TypeChecker) {
  if (sym) {
    const typ = checker.getTypeOfSymbolAtLocation(sym, node)
    if (typ.aliasTypeArguments) {
      let index = 0
      if (checker.typeToString(typ).startsWith('TSaveActionEvent')) {
        index = 1
      }
      return convertType(typ.aliasTypeArguments[index], checker)
    }
  }
}

function serializeBody(method: ts.MethodDeclaration, checker: ts.TypeChecker) {
  const symbol = checker.getSymbolAtLocation(method.name)
  let data
  let error
  if (symbol) {
    const returnTypeNode = method.type!

    let dataIndex = 0
    let errorIndex = 1
    const typeName = checker.typeToString(checker.getTypeFromTypeNode(returnTypeNode))
    if (typeName.startsWith('Promise<TSaveStatePayload')) {
      dataIndex = 1
      errorIndex = 2
    }
    if (returnTypeNode && ts.isTypeReferenceNode(returnTypeNode)) {
      const typeData = checker.getTypeFromTypeNode(returnTypeNode.typeArguments![dataIndex])
      data = convertType(typeData, checker)
      const errorArgument = returnTypeNode.typeArguments![errorIndex]
      if (errorArgument) {
        const typeError = checker.getTypeFromTypeNode(errorArgument)
        error = convertType(typeError, checker)
      }
    }
  }
  return { data, error }
}

function getIntentAuthType(sym: ts.Symbol | undefined, node: ts.Node, checker: ts.TypeChecker) {
  if (sym) {
    const typ = checker.getTypeOfSymbolAtLocation(sym, node)
    if (typ.aliasTypeArguments) {
      let index = 1
      if (checker.typeToString(typ).startsWith('TSaveActionEvent')) {
        index = 2
      }
      const authType = checker.typeToString(typ.aliasTypeArguments[index])
      switch (authType) {
        case 'TBaseAuthContext<{ apiKey: string; }>':
          return 'APIKEY'
        case 'TBaseAuthContext<{ username: string; password: string; }>':
          return 'BASIC'
        case 'TBaseAuthContext<{ accessToken: string; }>':
          return 'OAUTH2'
        case 'TBaseAuthContext<undefined>':
          return 'NONE'
        default:
          return authType
      }
    }
  }
}
export function intentTypesToSchemaConverter(intentPath: string, options: ts.CompilerOptions): TOutput {
  const program = ts.createProgram([intentPath], options)
  const sourceFile = program.getSourceFile(intentPath)
  const checker = program.getTypeChecker()

  const output: TOutput = {
    requestBody: {},
    response: {}
  }

  if (sourceFile) {
    ts.forEachChild(sourceFile, visit)
  }

  return output

  function visit(node: ts.Node) {
    // Only consider exported nodes
    if (!isNodeExported(node)) {
      return
    }

    if (ts.isClassDeclaration(node) && node.name) {
      const actionMethodNode = findActionMethod(node)
      if (actionMethodNode) {
        // This is a top level class, get its symbol

        const parameterNode = actionMethodNode.parameters[0]
        const sym = checker.getSymbolAtLocation(parameterNode.name)
        output.requestBody = serializeParameters(sym, parameterNode, checker)
        output.response = serializeBody(actionMethodNode, checker)
        output.intentAuthType = getIntentAuthType(sym, parameterNode, checker)
      }
    }
  }

  function isNodeExported(node: ts.Node): boolean {
    return (
      (ts.getCombinedModifierFlags(node as ts.Declaration) & ts.ModifierFlags.Export) !== 0 ||
      (!!node.parent && node.parent.kind === ts.SyntaxKind.SourceFile)
    )
  }
}

export default function generator({
  intentsDir,
  intents,
  integrationUuid,
  integrationName
}: {
  intentsDir: string
  intents: string[]
  integrationUuid: string
  integrationName: string
}): string {
  const template = fs.readFileSync('./src/templates/openapi.json.mustache', 'utf8')
  const schemas = intents.map(intent => {
    const intentPath = path.join(intentsDir, `${intent}.ts`)
    const typeSchema = intentTypesToSchemaConverter(intentPath, {
      target: ts.ScriptTarget.ES5,
      module: ts.ModuleKind.CommonJS
    })
    return {
      intentName: intent,
      typeSchema: {
        response: JSON.stringify(typeSchema.response),
        requestBody: JSON.stringify(typeSchema.requestBody)
      },
      oauth2: typeSchema.intentAuthType === 'OAUTH2'
    }
  })
  return Mustache.render(template, {
    intents,
    integrationUuid,
    integrationName,
    schemas
  })
}
