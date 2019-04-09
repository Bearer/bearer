import * as functions from '@bearer/functions'
import * as fs from 'fs'
import * as globby from 'globby'
import * as path from 'path'
import * as ts from 'typescript'
import * as TJS from 'typescript-json-schema'
import specGenerator from '@bearer/openapi-generator'

const config = ts.readConfigFile(path.join(__dirname, '../../templates/start', 'tsconfig.json'), ts.sys.readFile)

const NON_FUNCTION_NAMES = ['DBClient']
const FUNCTION_NAMES = Object.keys(functions).filter(functionName => !NON_FUNCTION_NAMES.includes(functionName))
const FUNCTION_TYPE_IDENTIFIER = 'functionType'

const functionEntries: IFunctionEntry[] = []

function bodySchema(tsType: ts.ClassDeclaration, generator: TJS.JsonSchemaGenerator): TJS.Definition {
  return {}
}

class FunctionNodeAdapter implements IFunctionEntry {
  constructor(
    readonly functionName: string,
    private readonly node: ts.ClassDeclaration,
    // @ts-ignore
    private readonly generator: TJS.JsonSchemaGenerator
  ) {}

  get functionClassName() {
    return getIdentifier(this.node).escapedText.toString()
  }

  get functionType() {
    return getPropertyValue(this.node, FUNCTION_TYPE_IDENTIFIER)
  }

  get paramsSchema() {
    return [...DEFAULT_PARAMS]
  }

  adaptParamsSchema({ properties = {} }: { properties?: { [key: string]: any } }): ISchemaParam[] {
    if (properties) {
      return Object.keys(properties).map(propName => {
        const { schema, required, name, description } = properties[propName]
        return {
          schema: schema || { type: 'string' },
          in: 'query',
          description: description || propName,
          required: required !== undefined ? required : true,
          name: name || propName
        } as ISchemaParam
      })
    }
    return []
  }

  get bodySchema() {
    return {
      properties: bodySchema(this.node, this.generator)
    }
  }

  get outputSchema() {
    return {}
  }

  get adapt(): IFunctionEntry {
    return {
      functionClassName: this.functionClassName,
      functionName: this.functionName,
      functionType: this.functionType,
      paramsSchema: this.paramsSchema,
      bodySchema: this.bodySchema,
      outputSchema: this.outputSchema
    }
  }
}

export function isFunctionClass(tsNode: ts.Node): boolean {
  if (!ts.isClassDeclaration(tsNode)) {
    return false
  }
  return extendsFunctionType(tsNode) && implementsAction(tsNode)
}

function extendsFunctionType(tsClass: ts.ClassDeclaration): boolean {
  const extendedClasses = (tsClass.heritageClauses || []).filter(hc => hc.token === ts.SyntaxKind.ExtendsKeyword)
  return Boolean(
    extendedClasses.find(hc =>
      Boolean(hc.types.find(t => FUNCTION_NAMES.includes((t.expression as ts.Identifier).escapedText.toString())))
    )
  )
}

function implementsAction(tsClass: ts.ClassDeclaration): boolean {
  return Boolean(
    tsClass.members.filter(m => {
      return ts.isMethodDeclaration(m)
    })
  )
}

export function getIdentifier(tsNode: ts.ClassDeclaration | ts.PropertyDeclaration): ts.Identifier {
  return tsNode.name as ts.Identifier
}

export function getFunctionName(tsSourceFile: ts.SourceFile): string {
  return path.basename(tsSourceFile.fileName).split('.')[0]
}

export function getPropertyValue(tsNode: ts.ClassDeclaration, propertyName: string): any {
  if (tsNode.members) {
    const declaration = tsNode.members.find(node => {
      return ts.isPropertyDeclaration(node) && (node.name as ts.Identifier).escapedText.toString() === propertyName
    }) as ts.PropertyDeclaration
    if (declaration && declaration.initializer) {
      return (declaration.initializer as ts.StringLiteral).text
    }
  }
}

export function transformer(generator: TJS.JsonSchemaGenerator): ts.TransformerFactory<ts.SourceFile> {
  return (context: ts.TransformationContext) => {
    return (tsSourceFile: ts.SourceFile) => {
      function visit(tsNode: ts.Node) {
        if (isFunctionClass(tsNode)) {
          const adapter = new FunctionNodeAdapter(
            getFunctionName(tsSourceFile),
            tsNode as ts.ClassDeclaration,
            generator
          )
          functionEntries.push(adapter.adapt)
        }
        return tsNode
      }
      return ts.visitEachChild(tsSourceFile, visit, context)
    }
  }
}

export class FunctionCodeProcessor {
  constructor(private readonly srcFunctionsDir: string, private readonly transformer: any) {}

  async run() {
    const files = await globby(`${this.srcFunctionsDir}/*.ts`)

    files.forEach(file => {
      const sourceFile = ts.createSourceFile(
        file,
        fs.readFileSync(file, 'utf8'),
        ts.ScriptTarget.Latest,
        false,
        ts.ScriptKind.TSX
      )
      ts.transform(sourceFile, [this.transformer])
    })
  }
}

export class OpenApiSpecGenerator {
  constructor(
    private readonly srcFunctionsDir: string,
    private readonly bearerConfig: { integrationTitle: string | undefined; buid: string }
  ) {}

  async build() {
    const files = await globby(`${this.srcFunctionsDir}/*.ts`)
    const programGenerator = TJS.getProgramFromFiles(
      files,
      {
        ...config.config.compilerOptions,
        allowUnusedLabels: true,
        // be indulgent
        noUnusedParameters: false,
        noUnusedLocals: false
      },
      './ok'
    )

    const generator = TJS.buildGenerator(programGenerator, {
      required: true
    })

    if (!generator) {
      throw new Error('Please fix above issues before')
    }

    files.forEach(file => {
      const sourceFile = ts.createSourceFile(
        file,
        fs.readFileSync(file, 'utf8'),
        ts.ScriptTarget.Latest,
        false,
        ts.ScriptKind.TSX
      )
      ts.transform(sourceFile, [transformer(generator)])
    })
    return specGenerator({
      functions: functionEntries.map(entry => entry.functionName),
      functionsDir: this.srcFunctionsDir,
      buid: this.bearerConfig.buid,
      integrationName: this.bearerConfig.integrationTitle || ''
    })
  }
}

const DEFAULT_PARAMS: ISchemaParam[] = [
  {
    in: 'header',
    name: 'authorization',
    schema: {
      type: 'string'
    },
    description: 'API Key',
    required: true
  }
]

export interface ISchemaParam {
  in: 'query' | 'header'
  schema: {
    type: 'string'
    format?: 'uuid'
  }
  description: string
  required: boolean
  name: string
}

export interface IOpenApiSpec {
  openapi: string
  info: any
  servers: any[]
  tags: any[]
  paths: {
    [key: string]: {
      post: {
        parameters: any[]
        summary: any[]
        requestBody: {
          content: {
            'application/json': {
              schema: any
            }
          }
        }
        responses: Record<TStatus, TResponse>
      }
    }
  }
}

type TStatus = '200' | '401' | '403'
type TResponse = {
  description: string
  content: {
    'application/json': {
      schema: any
    }
  }
}

interface IFunctionEntry {
  functionClassName: string
  functionType: string
  functionName: string
  paramsSchema: ISchemaParam[]
  bodySchema: any
  outputSchema: any
}
