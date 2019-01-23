import * as intents from '@bearer/intents'
import * as fs from 'fs'
import * as globby from 'globby'
import * as path from 'path'
import * as ts from 'typescript'
import * as TJS from 'typescript-json-schema'

const config = ts.readConfigFile(path.join(__dirname, '../../templates/start', 'tsconfig.json'), ts.sys.readFile)

const NON_INTENT_NAMES = ['DBClient']
const INTENT_NAMES = Object.keys(intents).filter(intentName => !NON_INTENT_NAMES.includes(intentName))
const INTENT_TYPE_IDENTIFIER = 'intentType'

const intentEntries: IIntentEntry[] = []

function bodySchema(tsType: ts.ClassDeclaration, generator: TJS.JsonSchemaGenerator): TJS.Definition {
  return {}
}

class IntentNodeAdapter implements IIntentEntry {
  constructor(
    readonly intentName: string,
    private readonly node: ts.ClassDeclaration,
    // @ts-ignore
    private readonly generator: TJS.JsonSchemaGenerator
  ) {}

  get intentClassName() {
    return getIdentifier(this.node).escapedText.toString()
  }

  get intentType() {
    return getPropertyValue(this.node, INTENT_TYPE_IDENTIFIER)
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
  if (!ts.isClassDeclaration(tsNode)) {
    return false
  }
  return extendsIntentType(tsNode) && implementsAction(tsNode)
}

function extendsIntentType(tsClass: ts.ClassDeclaration): boolean {
  const extendedClasses = (tsClass.heritageClauses || []).filter(hc => hc.token === ts.SyntaxKind.ExtendsKeyword)
  return Boolean(
    extendedClasses.find(hc =>
      Boolean(hc.types.find(t => INTENT_NAMES.includes((t.expression as ts.Identifier).escapedText.toString())))
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

export function getIntentName(tsSourceFile: ts.SourceFile): string {
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
        if (isIntentClass(tsNode)) {
          const adapter = new IntentNodeAdapter(getIntentName(tsSourceFile), tsNode as ts.ClassDeclaration, generator)
          intentEntries.push(adapter.adapt)
        }
        return tsNode
      }
      return ts.visitEachChild(tsSourceFile, visit, context)
    }
  }
}

export class IntentCodeProcessor {
  constructor(private readonly srcIntentsDir: string, private readonly transformer: any) {}

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
  ) {}

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

  generate(
    entries: IIntentEntry[],
    { scenarioTitle, scenarioUuid }: { scenarioTitle?: string; scenarioUuid: string }
  ): IOpenApiSpec {
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
                ...unAuthorizedReponse,
                ...forbiddenResponse,
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
  // we probably don't neet this one anymore
  // ,{
  //   name: 'authId',
  //   schema: {
  //     type: 'string',
  //     format: 'uuid'
  //   },
  //   in: 'query',
  //   description: 'User Identifier',
  //   required: true
  // }
]

const unAuthorizedReponse = {
  '401': {
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
  }
}

const forbiddenResponse = {
  '403': {
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
  }
}

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

interface IIntentEntry {
  intentClassName: string
  intentType: string
  intentName: string
  paramsSchema: ISchemaParam[]
  bodySchema: any
  outputSchema: any
}
