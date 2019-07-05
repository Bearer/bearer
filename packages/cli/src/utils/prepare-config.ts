import * as fs from 'fs-extra'

import * as ts from 'typescript'

import { getFunctionName, FunctionCodeProcessor, isFunctionClass } from './generators'

type TConfig = {
  functions: string[]
  buid: string
  auth?: any
}

export const transformer = (functions: string[]) => (context: ts.TransformationContext) => {
  return (tsSourceFile: ts.SourceFile) => {
    function visit(tsNode: ts.Node) {
      if (isFunctionClass(tsNode)) {
        const functionName = getFunctionName(tsSourceFile)
        functions.push(functionName)
      }
      return tsNode
    }
    return ts.visitEachChild(tsSourceFile, visit, context)
  }
}

export default (authConfigFile: string, buid: string, functionsDir: string): Promise<TConfig> => {
  return new Promise((resolve, reject) => {
    const functions: string[] = []
    new FunctionCodeProcessor(functionsDir, transformer(functions))
      .run()
      .then(() => {
        const content = fs.readFileSync(authConfigFile, { encoding: 'utf8' })
        const config: TConfig = { functions, buid, auth: JSON.parse(content) }
        resolve(config)
      })
      .catch(error => {
        reject(error)
      })
  })
}

export const HANDLER_NAME = 'index'
export const HANDLER_NAME_WITH_EXT = [HANDLER_NAME, 'js'].join('.')
