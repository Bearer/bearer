import * as ts from 'typescript'
import * as fs from 'fs-extra'
import * as path from 'path'

import { getSourceCode } from './utils'
import ReplaceIntentDecorators from './transformers/replace-intent-decorator'
import BearerScenarioIdInjector from './transformers/scenario-id-accessor-injector'
import PropInjector from './transformers/prop-injector'
import PropBearerContextInjector from './transformers/prop-bearer-context-injector'
import PropImporter from './transformers/prop-importer'

export default class Transpiler {
  private service: ts.LanguageService
  private rootFileNames: string[] = []
  private subscribers: ts.MapLike<Array<() => void>> = {}

  constructor(
    private ROOT_DIRECTORY = process.cwd(),
    private watchFiles = true
  ) {
    const config = ts.readConfigFile(
      path.join(this.BUILD_DIRECTORY, 'tsconfig.json'),
      ts.sys.readFile
    )
    const parsed = ts.parseJsonConfigFileContent(
      config,
      ts.sys,
      this.SRC_DIRECTORY
    )
    this.rootFileNames = parsed.fileNames
  }

  run(
    options: ts.CompilerOptions = {
      module: ts.ModuleKind.CommonJS
    }
  ) {
    // ensure global directory: quick and dirty
    // TODO: find another way to have the global present within src directory
    const files: ts.MapLike<{ version: number }> = {}
    const servicesHost: ts.LanguageServiceHost = {
      getScriptFileNames: () => this.rootFileNames,
      getScriptVersion: fileName =>
        files[fileName] && files[fileName].version.toString(),
      getScriptSnapshot: fileName => {
        if (!fs.existsSync(fileName)) {
          return undefined
        }
        return ts.ScriptSnapshot.fromString(
          fs.readFileSync(fileName).toString()
        )
      },
      getCurrentDirectory: () => process.cwd(),
      getCompilationSettings: () => options,
      getDefaultLibFileName: options => ts.getDefaultLibFilePath(options),
      getCustomTransformers: () => this.transformers,
      fileExists: ts.sys.fileExists,
      readFile: ts.sys.readFile,
      readDirectory: ts.sys.readDirectory
    }
    // Create the language service files
    this.service = ts.createLanguageService(
      servicesHost,
      ts.createDocumentRegistry()
    )
    // Now let's watch the files
    this.rootFileNames.forEach(fileName => {
      files[fileName] = { version: 0 }
      // First time around, emit all files
      this.emitFile(fileName)

      if (this.watchFiles) {
        // Add a watch on the file to handle next change
        fs.watchFile(
          fileName,
          { persistent: true, interval: 250 },
          (curr, prev) => {
            // Check timestamp
            if (+curr.mtime <= +prev.mtime) {
              return
            }
            // Update the version to signal a change in the file
            files[fileName].version++
            // write the changes to disk
            this.emitFile(fileName)
          }
        )
      }
    })
    if (!this.watchFiles) {
      this.stop()
      process.exit(0)
    }
  }

  stop() {
    this.rootFileNames.forEach(fileName => {
      fs.unwatchFile(fileName)
    })
    this.trigger('STOP')
  }

  on(event: string, callback: () => void) {
    this.subscribers[event] = this.subscribers[event] || []
    this.subscribers[event].push(callback)
  }

  private trigger = (eventName: string) => {
    ;(this.subscribers[eventName] || []).forEach(callback => callback())
  }

  get transformers(): ts.CustomTransformers {
    const verbose = true
    return {
      before: [
        ReplaceIntentDecorators({ verbose }),
        BearerScenarioIdInjector({ verbose }),
        PropImporter({ verbose }),
        PropInjector({ verbose }),
        PropBearerContextInjector({ verbose }),
        dumpSourceCode(this.SRC_DIRECTORY, this.BUILD_SRC_DIRECTORY)({
          verbose: true
        })
      ],
      after: []
    }
  }

  emitFile(fileName: string) {
    const output = this.service.getEmitOutput(fileName)

    if (!output.emitSkipped) {
      console.log(`Emitting ${fileName}`)
    } else {
      console.log(`Emitting ${fileName} failed`)
      this.logErrors(fileName)
    }
  }

  logErrors(fileName: string) {
    let allDiagnostics = this.service
      .getCompilerOptionsDiagnostics()
      .concat(this.service.getSyntacticDiagnostics(fileName))
      .concat(this.service.getSemanticDiagnostics(fileName))

    allDiagnostics.forEach(diagnostic => {
      let message = ts.flattenDiagnosticMessageText(
        diagnostic.messageText,
        '\n'
      )
      if (diagnostic.file) {
        let { line, character } = diagnostic.file.getLineAndCharacterOfPosition(
          diagnostic.start!
        )
        console.log(
          `  Error ${diagnostic.file.fileName} (${line + 1},${character +
            1}): ${message}`
        )
      } else {
        console.log(`  Error: ${message}`)
      }
    })
  }

  private get SRC_DIRECTORY(): string {
    return path.join(this.ROOT_DIRECTORY, 'screens')
  }

  private get BUILD_DIRECTORY(): string {
    return path.join(this.ROOT_DIRECTORY, '.build')
  }

  private get BUILD_SRC_DIRECTORY(): string {
    return path.join(this.BUILD_DIRECTORY, 'src')
  }
}

type TransformerOptions = {
  verbose?: true
}

function dumpSourceCode(srcDirectory, buildDirectory) {
  return function storeOutput({
    verbose
  }: TransformerOptions = {}): ts.TransformerFactory<ts.SourceFile> {
    return transformContext => {
      return tsSourceFile => {
        let outPath = tsSourceFile.fileName
          .replace(srcDirectory, buildDirectory)
          .replace(/js$/, 'ts')
          .replace(/jsx$/, 'tsx')

        fs.ensureFileSync(outPath)
        fs.writeFileSync(outPath, getSourceCode(tsSourceFile))

        return tsSourceFile
      }
    }
  }
}
