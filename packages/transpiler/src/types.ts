import { TInputDecoratorOptions, TOutputDecoratorOptions } from '@bearer/types/lib/input-output-decorators'
import * as ts from 'typescript'

import Metadata from './metadata'

export type ComponentMetadata = {
  classname: string
  isRoot: boolean
  initialTagName: string
  fileName: string
  finalTagName: string
  name?: string
  imports?: string[]
  inputs?: TComponentInputDefinition[]
  outputs?: TComponentOutputDefinition[]
}

export type TComponentInputDefinition = {
  name: string
  type: 'string' | 'number'
  default: string | number
}

type TBasicFormat = 'string' | 'number' | 'boolean' | 'object' | 'any'
type TUnknown = 'unspecified'
export type TOuputFormat = TUnknown | { type: TBasicFormat } | { [key: string]: TBasicFormat | TBasicFormat[] }

export type TComponentOutputDefinition = {
  name: string
  eventName: string
  payloadFormat: TOuputFormat
}

export type TransformerOptions = {
  verbose?: true
  metadata?: Metadata
}

export type FileTransformerOptions = TransformerOptions & {
  outDir: string
  srcDir?: string
}

export type SourceCodeTransformerOptions = TransformerOptions & {
  srcDirectory: string
  buildDirectory: string
}

export type SpecComponent = {
  classname: string
  isRoot: boolean
  initialTagName: string
  group: string
  label: string
  input: any
  output: any
}

export type RootComponent = SpecComponent & {
  finalTagName: string
}

export type CompileSpec = { components: SpecComponent[] }

export type InputMeta = TInputDecoratorOptions & {
  propDeclarationName: string
  typeIdentifier?: ts.TypeNode
  intializer?: ts.Expression
  loadMethodName: string
  intentMethodName: string
  watcherName: string
}

export type CreateFetcherMeta = {
  intentName: string
  intentMethodName: string
}

export type TCreateLoadResourceMethod = {
  propertyReferenceIdName: string
  typeIdentifier?: ts.TypeNode
  propDeclarationName: string
  intentMethodName: string
  intentReferenceIdKeyName: string
  loadMethodName: string
  intentArguments?: string[]
}

export type TCreateLoadDataCall = {
  loadMethodName: string
}

export type TAddAutoLoad = {
  propertyReferenceIdName: string
  autoLoad: boolean
  loadMethodName: string
}

export type OutputMeta = TOutputDecoratorOptions & {
  propDeclarationName: string
  propDeclarationNameRefId: string
  intentMethodName: string
  intentName: string
  loadMethodName: string
  initializer: ts.Expression
  typeIdentifier?: ts.TypeNode
  propertyReferenceIdName: string
  autoLoad?: true | false
}
