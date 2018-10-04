import Metadata from './metadata'

export type ComponentMetadata = {
  classname: string
  isRoot: boolean
  initialTagName: string
  finalTagName: string
  group?: string
  imports?: Array<string>
  inputs?: Array<TComponentInputDefinition>
  outputs?: Array<TComponentInputDefinition>
}

type TComponentInputDefinition = {
  name: string
  type: 'string' | 'number'
  default: string | number
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

export type CompileSpec = { components: Array<SpecComponent> }
