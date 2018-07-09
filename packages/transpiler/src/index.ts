import * as ts from 'typescript'
import * as fs from 'fs-extra'
import * as path from 'path'
import * as chokidar from 'chokidar'
import { getSourceCode } from './utils'
import PropInjector from './transformers/prop-injector'
import PropImporter from './transformers/prop-importer'

export default class Transpiler {
  private watcher: any
  private service: ts.LanguageService
  private rootFileNames: string[] = []
  constructor(
    private readonly SCREENS_DIRECTORY = path.join(process.cwd(), 'screens')
  ) {
    const config = ts.readConfigFile(
      path.join(this.BUILD_DIRECTORY, 'tsconfig.json'),
      ts.sys.readFile
    )
    const parsed = ts.parseJsonConfigFileContent(config, ts.sys, process.cwd())
    this.rootFileNames = parsed.fileNames
  }

  run(
    options: ts.CompilerOptions = {
      module: ts.ModuleKind.CommonJS
    }
  ) {
    this.watchNonTSFiles()

    fs.emptyDirSync(this.BUILD_SRC_DIRECTORY)

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
    })
  }

  stop() {
    if (this.watcher) {
      this.watcher.close()
    }
  }

  get transformers(): ts.CustomTransformers {
    return {
      before: [
        PropImporter({ verbose: true }),
        PropInjector({ verbose: true }),
        dumpSourceCode(this.SCREENS_DIRECTORY, this.BUILD_DIRECTORY)({
          verbose: true
        })
      ],
      after: []
    }
  }

  private watchNonTSFiles() {
    function callback(error) {
      if (error) {
        console.log('error', error)
      }
    }

    this.watcher = chokidar.watch(this.SRC_DIRECTORY + '/**', {
      ignored: /\.tsx?$/,
      persistent: true,
      followSymlinks: false
    })

    this.watcher.on('all', (event, filePath) => {
      const relativePath = filePath.replace(this.SRC_DIRECTORY, '')
      const targetPath = path.join(this.BUILD_SRC_DIRECTORY, relativePath)
      // Creating symlink
      if (event == 'add' || event == 'addDir') {
        console.log('creating symlink')
        fs.ensureSymlink(filePath, targetPath, callback)
      }

      // // Deleting symlink
      if (event == 'unlink') {
        console.log('deleting symlink')
        fs.unlink(targetPath, err => {
          if (err) throw err
          console.log(targetPath + ' was deleted')
        })
      }
    })
  }

  emitFile(fileName: string) {
    let output = this.service.getEmitOutput(fileName)

    if (!output.emitSkipped) {
      console.log(`Emitting ${fileName} like a pro!!!`)
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
    return path.join(this.SCREENS_DIRECTORY, 'src')
  }

  private get BUILD_DIRECTORY(): string {
    return path.join(this.SCREENS_DIRECTORY, '.build')
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
