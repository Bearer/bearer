import * as ts from 'typescript'
import * as path from 'path'
import * as fs from 'fs-extra'

export interface TranspilerConfiguration extends ts.CompilerOptions {
  outDir: string
}
export default class Transpiler {
  options: TranspilerConfiguration
  constructor(options: TranspilerConfiguration) {
    this.options = options
  }
  async run(fileNames) {
    const { outDir } = this.options
    await fs.mkdirp(outDir)
    let program = ts.createProgram(fileNames, this.options)
    program.getRootFileNames().forEach(filePath => {
      let fileName = path.basename(filePath)
      fs.writeFileSync(path.join(outDir, fileName), fs.readFileSync(filePath))
    })
  }
}
