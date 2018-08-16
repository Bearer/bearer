import * as fs from 'fs-extra'
import * as path from 'path'
import * as ts from 'typescript'

import BearerStateInjector from './transformers/bearer-state-injector'
import ComponenttagNameScoping from './transformers/component-tag-name-scoping'
import GatherMetadata from './transformers/gather-metadata'
import generateMetadataFile from './transformers/generate-metadata-file'
import ImportsImporter from './transformers/imports-transformer'
import NavigatorScreenTransformer from './transformers/navigator-screen-transformer'
import PropBearerContextInjector from './transformers/prop-bearer-context-injector'
import PropImporter from './transformers/prop-importer'
import PropInjector from './transformers/prop-injector'
import BearerReferenceIdInjector from './transformers/reference-id-injector'
import ReplaceIntentDecorators from './transformers/replace-intent-decorator'
import RootComponentTransformer from './transformers/root-component-transformer'
import BearerScenarioIdInjector from './transformers/scenario-id-accessor-injector'
import { Metadata, SourceCodeTransformerOptions } from './types'
import { getSourceCode } from './utils'

export type TranpilerOptions = {
  ROOT_DIRECTORY?: string
  watchFiles?: boolean
  buildFolder?: string
  srcFolder?: string
  verbose?: boolean
  tagNamePrefix?: string
  tagNameSuffix?: string
}

export default class Transpiler {
  get transformers(): ts.CustomTransformers {
    const verbose = true
    return {
      before: [
        GatherMetadata({ verbose, metadata: this.metadata }),
        RootComponentTransformer({ verbose, metadata: this.metadata }),
        BearerReferenceIdInjector({ verbose, metadata: this.metadata }),
        ReplaceIntentDecorators({ verbose, metadata: this.metadata }),
        BearerScenarioIdInjector({ verbose, metadata: this.metadata }),
        PropImporter({ verbose, metadata: this.metadata }),
        PropInjector({ verbose, metadata: this.metadata }),
        PropBearerContextInjector({ verbose, metadata: this.metadata }),
        BearerStateInjector({ verbose, metadata: this.metadata }),
        NavigatorScreenTransformer({ verbose, metadata: this.metadata }),
        ImportsImporter({ verbose, metadata: this.metadata }),
        ComponenttagNameScoping({ verbose, metadata: this.metadata }),
        dumpSourceCode({
          verbose: true,
          srcDirectory: this.VIEWS_DIRECTORY,
          buildDirectory: this.BUILD_SCR_DIRECTORY
        }),
        generateMetadataFile({ verbose, metadata: this.metadata, outDir: this.BUILD_SCR_DIRECTORY })
      ],
      after: []
    }
  }

  private get BUILD_DIRECTORY(): string {
    return path.join(this.ROOT_DIRECTORY, this.buildFolder)
  }

  private get BUILD_SCR_DIRECTORY(): string {
    return path.join(this.BUILD_DIRECTORY, 'src')
  }

  private get VIEWS_DIRECTORY(): string {
    return path.join(this.ROOT_DIRECTORY, this.srcFolder)
  }
  private service: ts.LanguageService
  private rootFileNames: string[] = []
  private subscribers: ts.MapLike<Array<() => void>> = {}

  private readonly ROOT_DIRECTORY
  private watchFiles = true
  private buildFolder = '.bearer/views'
  private srcFolder = 'views'
  private verbose = true
  private files: ts.MapLike<{ version: number }> = {}
  private metadata: Metadata = {
    components: []
  }

  private compilerOptionsts: ts.CompilerOptions = {
    module: ts.ModuleKind.CommonJS
  }

  constructor(options?: Partial<TranpilerOptions>) {
    Object.assign(this, options)

    this.ROOT_DIRECTORY = this.ROOT_DIRECTORY || process.cwd()

    if (options.tagNamePrefix) {
      this.metadata.prefix = options.tagNamePrefix
    }

    if (options.tagNameSuffix) {
      this.metadata.suffix = options.tagNameSuffix
    }
  }

  run() {
    this.refresh()

    if (!this.watchFiles) {
      this.stop()
    }
  }

  emitFiles = () => {
    // Now let's watch the files
    this.rootFileNames.forEach(fileName => {
      this.files[fileName] = { version: 0 }
      // First time around, emit all files
      this.emitFile(fileName)
      if (this.watchFiles) {
        // Add a watch on the file to handle next change
        fs.watchFile(fileName, { persistent: true, interval: 250 }, (curr, prev) => {
          // Check timestamp
          if (+curr.mtime <= +prev.mtime) {
            return
          }
          // Update the version to signal a change in the file
          this.files[fileName].version++
          // write the changes to disk
          this.emitFile(fileName)
        })
      }
    })
  }

  refresh() {
    this.clearWatchers()

    const config = ts.readConfigFile(path.join(this.BUILD_DIRECTORY, 'tsconfig.json'), ts.sys.readFile)

    if (config.error) {
      throw new Error(config.error.messageText as string)
    }

    const parsed = ts.parseJsonConfigFileContent(config, ts.sys, this.VIEWS_DIRECTORY)
    this.rootFileNames = parsed.fileNames
    if (!this.rootFileNames.length) {
      console.warn('[BEARER]', 'No file to transpile')
    }

    const servicesHost: ts.LanguageServiceHost = {
      getScriptFileNames: () => this.rootFileNames,
      getScriptVersion: fileName => this.files[fileName] && this.files[fileName].version.toString(),
      getScriptSnapshot: fileName => {
        if (!fs.existsSync(fileName)) {
          return null
        }

        return ts.ScriptSnapshot.fromString(fs.readFileSync(fileName).toString())
      },
      getCurrentDirectory: () => process.cwd(),
      getCompilationSettings: () => this.compilerOptionsts,
      getDefaultLibFileName: options => ts.getDefaultLibFilePath(options),
      getCustomTransformers: () => this.transformers,
      fileExists: ts.sys.fileExists,
      readFile: ts.sys.readFile,
      readDirectory: ts.sys.readDirectory
    }
    // Create the language service files
    this.service = ts.createLanguageService(servicesHost, ts.createDocumentRegistry())

    this.emitFiles()
  }

  stop() {
    this.clearWatchers()
    this.trigger('STOP')
  }

  clearWatchers(): void {
    this.rootFileNames.forEach(fileName => {
      fs.unwatchFile(fileName)
    })
  }

  on(event: string, callback: () => void) {
    this.subscribers[event] = this.subscribers[event] || []
    this.subscribers[event].push(callback)
  }

  emitFile = (fileName: string) => {
    const output = this.service.getEmitOutput(fileName)

    if (!output.emitSkipped) {
      if (this.verbose) {
        console.log(`Emitting ${fileName}`)
      }
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
      let message = ts.flattenDiagnosticMessageText(diagnostic.messageText, '\n')
      if (diagnostic.file) {
        let { line, character } = diagnostic.file.getLineAndCharacterOfPosition(diagnostic.start!)
        console.log(`  Error ${diagnostic.file.fileName} (${line + 1},${character + 1}): ${message}`)
      } else {
        console.log(`  Error: ${message}`)
      }
    })
  }

  private trigger = (eventName: string) => {
    const subscribers = this.subscribers[eventName] || []
    subscribers.forEach(callback => {
      callback()
    })
  }
}

function dumpSourceCode(
  { srcDirectory, buildDirectory }: SourceCodeTransformerOptions = { srcDirectory, buildDirectory }
): ts.TransformerFactory<ts.SourceFile> {
  return _transformContext => {
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
