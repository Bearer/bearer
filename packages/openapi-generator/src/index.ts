import path from 'path'
import ts from 'typescript'
import { topOfSpec, specPath } from './openapi-template'

import { convertType } from './convert-type'

export const INTENT_ACTION = 'action'
type TOutput = { requestBody: any; response: any; intentAuthType?: string }

/**
 * find action method in class and return relevant node
 */
function findActionMethod(node: ts.ClassDeclaration): ts.MethodDeclaration | ts.ArrowFunction | undefined {
  const action = node.members.find((node: ts.ClassElement) => {
    if (node.name) {
      return node.name.getText() === INTENT_ACTION
    }
    return false
  })
  if (ts.isPropertyDeclaration(action!)) {
    return (action! as ts.PropertyDeclaration).initializer as ts.ArrowFunction
  }
  return action as ts.MethodDeclaration
}

function serializeBody(method: ts.MethodDeclaration | ts.ArrowFunction, checker: ts.TypeChecker) {
  const symbol = checker.getSymbolAtLocation(method.name)
  let data
  let error
  if (symbol) {
    const returnTypeNode = method.type!

    let dataIndex = 0 // index of ReturnedData in TFetchPromise<ReturnedData, ReturnedError>
    let errorIndex = 1 // index of ReturnedError in TFetchPromise<ReturnedData, ReturnedError>
    const typeName = checker.typeToString(checker.getTypeFromTypeNode(returnTypeNode))
    if (typeName.startsWith('Promise<TSaveStatePayload')) {
      dataIndex = 1 // index of ReturnedData in TSavePromise<State, ReturnedData, ReturnedError>
      errorIndex = 2 // index of ReturnedData in TSavePromise<State, ReturnedData, ReturnedError>
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

function resolveParameterTypeIndex(t: ts.Type, checker: ts.TypeChecker, i: number, j: number): number {
  let index = i
  if (checker.typeToString(t).startsWith('TSaveActionEvent')) {
    index = j
  }
  return index
}
/**
 *  Find the Parameter type in intent and convert it to json-schema
 */
function serializeParameters(sym: ts.Symbol | undefined, node: ts.Node, checker: ts.TypeChecker) {
  if (sym) {
    const typ = checker.getTypeOfSymbolAtLocation(sym, node)
    if (typ.aliasTypeArguments) {
      const index = resolveParameterTypeIndex(typ, checker, 0, 1)
      return convertType(typ.aliasTypeArguments[index!], checker)
    }
  }
}

/**
 *  Find the auth context type and return it's name
 *  @param sym action method symbol
 *  @param node action method symbol value declaration
 *  @param checker program type checker
 */
function getIntentAuthType(sym: ts.Symbol | undefined, node: ts.Node, checker: ts.TypeChecker): string | undefined {
  if (sym) {
    const typ = checker.getTypeOfSymbolAtLocation(sym, node)
    if (typ.aliasTypeArguments) {
      const index = resolveParameterTypeIndex(typ, checker, 1, 2)
      const authType = checker.typeToString(typ.aliasTypeArguments[index])
      if (/apiKey/.test(authType)) return 'APIKEY'
      if (/username/.test(authType) && /password/.test(authType)) return 'BASIC'
      if (/accessToken/.test(authType)) return 'OAUTH2'
      if (/undefined/.test(authType)) return 'NONE'
      return authType
    }
  }
  return
}

/**
 *  Convert single intent types to json schema
 *  This function tries to find the relevant type definitions, Params and `action` method resposne type
 *  and converts the types to json-schema
 *  @param intentPath absolute path to intent
 *  @param options compiler options
 */
export function intentTypesToSchemaConverter(
  intentPath: string,
  options: ts.CompilerOptions = { target: ts.ScriptTarget.ES5, module: ts.ModuleKind.CommonJS }
): TOutput {
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

/**
 *  Generate full openapi spec from intents
 *  @param intentsDir absolute path to intents directory
 *  @param intents list of intent names
 *  @param integrationUuid integration unique identifier
 *  @param integrationName name of the integration
 */
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
  const doc = topOfSpec(integrationName)
  const schemas = intents.reduce((acc, intent) => {
    const intentPath = path.join(intentsDir, `${intent}.ts`)
    const typeSchema = intentTypesToSchemaConverter(intentPath)
    return Object.assign(
      acc,
      specPath({
        integrationUuid,
        intentName: intent,
        response: { type: 'object', properties: typeSchema.response },
        requestBody: typeSchema.requestBody,
        oauth2: typeSchema.intentAuthType === 'OAUTH2'
      })
    )
  }, {})

  return JSON.stringify(Object.assign(doc, { paths: schemas }), null, 2)
}
