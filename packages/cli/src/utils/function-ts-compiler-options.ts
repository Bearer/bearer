import * as ts from 'typescript'
const options: ts.CompilerOptions = {
  allowUnreachableCode: false,
  declaration: false,
  lib: ['es2017'],
  noUnusedLocals: false,
  noUnusedParameters: false,
  allowSyntheticDefaultImports: true,
  experimentalDecorators: true,
  moduleResolution: ts.ModuleResolutionKind.NodeJs,
  module: ts.ModuleKind.CommonJS,
  target: ts.ScriptTarget.ES2017
}
export default options
