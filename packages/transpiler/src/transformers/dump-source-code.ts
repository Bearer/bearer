import * as fs from 'fs-extra'
import * as ts from 'typescript'

import { SourceCodeTransformerOptions } from '../types'
import { getSourceCode } from '../utils'
import debug from '../logger'
const logger = debug.extend('source-code-dumping')

export default function dumpSourceCode(
  { srcDirectory, buildDirectory }: SourceCodeTransformerOptions = {
    srcDirectory,
    buildDirectory
  }
): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
    return tsSourceFile => {
      logger('processing %s', tsSourceFile.fileName)
      const outPath = tsSourceFile.fileName
        .replace(srcDirectory, buildDirectory)
        .replace(/js$/, 'ts')
        .replace(/jsx$/, 'tsx')
      fs.ensureFileSync(outPath)
      fs.writeFileSync(outPath, getSourceCode(tsSourceFile))

      return tsSourceFile
    }
  }
}
