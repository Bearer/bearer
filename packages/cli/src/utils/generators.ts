import * as intents from '@bearer/intents'
import * as fs from 'fs-extra'
import * as globby from 'globby'
import * as path from 'path'
import * as ts from 'typescript'
import * as TJS from 'typescript-json-schema'

const config = ts.readConfigFile(path.join(__dirname, '../../templates/start', 'tsconfig.json'), ts.sys.readFile)

const NON_INTENT_NAMES = ['DBClient']
const INTENT_NAMES = Object.keys(intents).filter(intentName => !NON_INTENT_NAMES.includes(intentName))
const INTENT_TYPE_IDENTIFIER = 'intentType'
const INTENT_NAME_IDENTIFIER = 'intentName'
const FUNCTION_NAME_IDENTIFIER = 'action'

let intentEntries: Array<IIntentEntry> = []

function initializerAsJson(tsType: ts.TypeReferenceNode | ts.TypeLiteralNode | ts.KeywordTypeNode, generator: any) {
  if (tsType) {
    switch (tsType.kind) {
      case ts.SyntaxKind.TypeReference:
        return generator.getSchemaForSymbol((tsType.typeName as ts.Identifier).escapedText.toString()) as any
      case ts.SyntaxKind.TypeLiteral:
      case ts.SyntaxKind.AnyKeyword:
        return { type: 'object' } as any
      case ts.SyntaxKind.NumberKeyword:
        return { type: 'number' }
      case ts.SyntaxKind.BooleanKeyword:
        return { type: 'boolean' }
      case ts.SyntaxKind.StringKeyword:
        return { type: 'string' }
      default:
        return {}
    }
  } else {
    return {}
  }
}

interface IIntentEntry {
  intentClassName: string
  intentType: string
  intentName: string
  paramsSchema: any
  bodySchema: any
  outputSchema: any
}

class IntentNodeAdapter implements IIntentEntry {
  constructor(private readonly node: ts.ClassDeclaration, private readonly generator: any) { }
  get intentClassName(): string {
    const identifier = getIdentifier(this.node)
    return identifier.escapedText.toString()
  }

  get intentType(): string {
    return getPropertyValue(this.node, INTENT_TYPE_IDENTIFIER)
  }

  get intentName(): string {
    return getPropertyValue(this.node, INTENT_NAME_IDENTIFIER)
  }

  get paramsSchema(): any {
    const typeNode = getFunctionParameterType(this.node, FUNCTION_NAME_IDENTIFIER, 'params') as ts.TypeReferenceNode

    const paramsSchema = initializerAsJson(typeNode, this.generator)
    return [...this.adaptParamsSchema(paramsSchema), ...this.defaultParams]
  }

  adaptParamsSchema(paramsSchema: any) {
    if (paramsSchema.properties) {
      return Object.keys(paramsSchema.properties).map(name => {
        return {
          in: 'query',
          schema: {
            type: 'string'
          },
          description: name,
          required: true,
          name
        }
      })
    }
    return []
  }

  get defaultParams(): any {
    return [
      {
        name: 'authorization',
        schema: {
          type: 'string'
        },
        in: 'header',
        description: 'API Key',
        required: true
      },
      {
        name: 'authId',
        schema: {
          type: 'string',
          format: 'uuid'
        },
        in: 'query',
        description: 'User Identifier',
        required: true
      }
    ]
  }

  get bodySchema(): any {
    const typeNode = getFunctionParameterType(this.node, FUNCTION_NAME_IDENTIFIER, 'body') as ts.TypeReferenceNode

    return initializerAsJson(typeNode, this.generator)
  }

  get outputSchema(): any {
    return {}
  }

  get adapt(): IIntentEntry {
    return {
      intentClassName: this.intentClassName,
      intentName: this.intentName,
      intentType: this.intentType,
      paramsSchema: this.paramsSchema,
      bodySchema: this.bodySchema,
      outputSchema: this.outputSchema
    }
  }
}

export function isIntentClass(tsNode: ts.Node): boolean {
  const isClass = ts.isClassDeclaration(tsNode) && tsNode.name

  const intentTypeValue = getPropertyValue(tsNode as ts.ClassDeclaration, INTENT_TYPE_IDENTIFIER)
  return (!!isClass) && INTENT_NAMES.includes(intentTypeValue)
}

export function getIdentifier(tsNode: ts.ClassDeclaration | ts.PropertyDeclaration): ts.Identifier {
  return tsNode.name as ts.Identifier
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

function getFunctionParameterType(
  tsNode: ts.ClassDeclaration,
  functionNameIdentifier: string,
  parameterIdentifier: string
) {
  if (tsNode.members) {
    const methodDeclaration = tsNode.members.find(node => {
      return (
        ts.isMethodDeclaration(node) && (node.name as ts.Identifier).escapedText.toString() === functionNameIdentifier
      )
    }) as ts.MethodDeclaration

    if (methodDeclaration.parameters) {
      const parameter = methodDeclaration.parameters.find(node => {
        return ts.isParameter(node) && (node.name as ts.Identifier).escapedText.toString() === parameterIdentifier
      }) as ts.ParameterDeclaration
      if (parameter) {
        return parameter.type
      }
    }
  }
  return null
}

export function transformer(generator: any) {
  return (context: ts.TransformationContext) => {
    function visit(tsNode: ts.Node) {
      if (isIntentClass(tsNode)) {
        const adapter = new IntentNodeAdapter(tsNode as ts.ClassDeclaration, generator)
        intentEntries.push(adapter.adapt)
      }
      return tsNode
    }
    return (tsSourceFile: ts.SourceFile) => {
      return ts.visitEachChild(tsSourceFile, visit, context)
    }
  }
}

export class IntentCodeProcessor {
  constructor(
    private readonly srcIntentsDir: string,
    private readonly transformer: any,
  ) { }

  async run() {
    const files = await globby(`${this.srcIntentsDir}/*.ts`)

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
    private readonly srcIntentsDir: string,
    private readonly bearerConfig: { scenarioTitle: string | undefined; scenarioUuid: string }
  ) { }

  async build() {
    const files = await globby(`${this.srcIntentsDir}/*.ts`)

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

    return this.generate(intentEntries, this.bearerConfig)
  }

  generate(entries: Array<IIntentEntry>, { scenarioTitle, scenarioUuid }: any): any {
    return {
      openapi: '3.0.0',
      info: {
        description: `openapi definition file for ${scenarioTitle}`,
        version: '0.0.1',
        title: scenarioTitle,
        contact: {
          email: 'bearer@bearer.sh'
        },
        license: {
          name: 'MIT'
        }
      },
      servers: [
        {
          url: 'https://int.bearer.sh/backend/api/v1'
        }
      ],
      tags: [
        {
          name: 'integration',
          description: `List of endpoints providing backend to backend integration with ${scenarioTitle}`,
          externalDocs: {
            description: 'Find out more',
            url: 'https://www.bearer.sh'
          }
        }
      ],
      paths: entries.reduce((acc, entry) => {
        return {
          ...acc,
          [`/${scenarioUuid}/${entry.intentName}`]: {
            post: {
              parameters: entry.paramsSchema,
              summary: entry.intentName,
              requestBody: {
                content: {
                  'application/json': {
                    schema: entry.bodySchema
                  }
                }
              },
              responses: {
                '401': {
                  description: 'Access forbidden',
                  content: {
                    'application/json': {
                      schema: {
                        type: 'object',
                        properties: {
                          error: {
                            type: 'string'
                          }
                        }
                      }
                    }
                  }
                },
                '403': {
                  description: 'Unauthorized',
                  content: {
                    'application/json': {
                      schema: {
                        type: 'object',
                        properties: {
                          error: {
                            type: 'string'
                          }
                        }
                      }
                    }
                  }
                },
                '200': {
                  description: 'Share',
                  content: {
                    'application/json': {
                      schema: entry.outputSchema
                    }
                  }
                }
              }
            }
          }
        }
      }, {})
    }
  }
}
