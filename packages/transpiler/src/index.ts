import * as ts from 'typescript'
import * as path from 'path'
import * as fs from 'fs-extra'
import { forEach } from 'async-foreach'
import ComponentTransformer from './component-transformer'

export interface TranspilerConfiguration extends ts.CompilerOptions {
  outDir: string
}

export interface Strategy {}
export default class Transpiler {
  constructor(private readonly options: TranspilerConfiguration) {}

  async run(fileNames) {
    const { outDir } = this.options
    await fs.mkdirp(outDir)
    let program = ts.createProgram(fileNames, this.options)
    let files = this.getFileNames(program)

    forEach(files, filePath => {
      let sourceFile = program.getSourceFile(filePath)
      let fileName = path.basename(filePath)
      program.getCompilerOptions()

      let res = ts.transform(sourceFile, [ComponentTransformer()])

      const resultFile = ts.createSourceFile(
        'tmp.ts',
        '',
        ts.ScriptTarget.Latest,
        /*setParentNodes*/ false,
        ts.ScriptKind.TS
      )
      const printer = ts.createPrinter({
        newLine: ts.NewLineKind.LineFeed
      })

      const result = printer.printNode(
        ts.EmitHint.Unspecified,
        res.transformed[0],
        resultFile
      )

      fs.writeFileSync(path.join(outDir, fileName), result)
    })
  }

  private getFileNames(program: ts.Program): string[] {
    return program
      .getSourceFiles()
      .filter(
        file =>
          !program.isSourceFileFromExternalLibrary(file) &&
          !file.fileName.endsWith('d.ts')
      )
      .map(file => file.fileName)
  }
}
