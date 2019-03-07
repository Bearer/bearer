import * as fs from 'fs-extra'
import * as path from 'path'
import * as ts from 'typescript'
import * as TJS from 'typescript-json-schema'

import debug from './logger'
const logger = debug.extend('compiler')
// optionally pass argument to schema generator
const settings: TJS.PartialArgs = {
  required: true,
  ignoreErrors: true
}

import BearerAuthorizedRequiredProp from './transformers/bearer-authorized-integration-id-prop-injector'
import bearerCleaning from './transformers/bearer-cleaning'
import BearerStateInjector from './transformers/bearer-state-injector'
import ComponenttagNameScoping from './transformers/component-tag-name-scoping'
import DumpSourceCode from './transformers/dump-source-code'
import EventNameScoping from './transformers/event-name-scoping'
import GatherIO from './transformers/gather-input-output'
import GatherMetadata from './transformers/gather-metadata'
import NavigatorScreenTransformer from './transformers/navigator-screen-transformer'
import PropBearerContextInjector from './transformers/prop-bearer-context-injector'
import BearerReferenceIdInjector from './transformers/reference-id-injector'
import ReplaceFunctionDecorators from './transformers/replace-function-decorator'
import RootComponentTransformer from './transformers/root-component-transformer'
import BearerIntegrationIdInjector from './transformers/integration-id-accessor-injector'
import EventNameNormalizer from './transformers/event-name-normalizer'
import I18nModifier from './transformers/i18n-modifier'

/*
 * Transformer modifying AST
 */
import InputDecoratorModifier from './transformers/input-decorator'
import OutputDecoratorModifier from './transformers/output-decorator'
import PropSetDecorator from './transformers/prop-set-decorator'

import { transformer as generateManifestFile } from './transformers/generate-manifest-file'

import Metadata from './metadata'

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
        InputDecoratorModifier({ verbose, metadata: this.metadata }),
        OutputDecoratorModifier({ verbose, metadata: this.metadata }),
        BearerReferenceIdInjector({ verbose, metadata: this.metadata }),
        ReplaceFunctionDecorators({ verbose, metadata: this.metadata }),
        BearerIntegrationIdInjector({ verbose, metadata: this.metadata }),
        PropBearerContextInjector({ verbose, metadata: this.metadata }),
        BearerStateInjector({ verbose, metadata: this.metadata }),
        NavigatorScreenTransformer({ verbose, metadata: this.metadata }),
        BearerAuthorizedRequiredProp({ verbose }),
        EventNameScoping({ metadata: this.metadata }),
        ComponenttagNameScoping({ verbose, metadata: this.metadata }),
        GatherIO({ verbose, metadata: this.metadata, generator: this.generator }),
        PropSetDecorator({ verbose, metadata: this.metadata }),
        bearerCleaning({ verbose, metadata: this.metadata }),
        EventNameNormalizer(),
        I18nModifier(),
        DumpSourceCode({
          verbose,
          srcDirectory: this.VIEWS_DIRECTORY,
          buildDirectory: this.BUILD_SRC_DIRECTORY
        })
      ],
      after: [
        generateManifestFile({
          metadata: this.metadata,
          outDir: this.BUILD_SRC_DIRECTORY,
          srcDir: this.ROOT_DIRECTORY
        })
      ]
    }
  }

  private get BUILD_DIRECTORY(): string {
    return path.join(this.ROOT_DIRECTORY, this.buildFolder)
  }

  private get BUILD_SRC_DIRECTORY(): string {
    return path.join(this.BUILD_DIRECTORY, 'src')
  }

  private get VIEWS_DIRECTORY(): string {
    return path.join(this.ROOT_DIRECTORY, this.srcFolder)
  }

  generator!: any
  private service: ts.LanguageService
  private rootFileNames: string[] = []
  private subscribers: ts.MapLike<(() => void)[]> = {}

  private readonly ROOT_DIRECTORY
  private watchFiles = true
  private buildFolder = '.bearer/views'
  private srcFolder = 'views'
  private verbose = true
  private files: ts.MapLike<{ version: number }> = {}
  private metadata: Metadata

  private compilerOptions: ts.CompilerOptions = {
    module: ts.ModuleKind.CommonJS
  }

  constructor(options?: Partial<TranpilerOptions>) {
    Object.assign(this, options)

    this.ROOT_DIRECTORY = this.ROOT_DIRECTORY || process.cwd()
    this.metadata = new Metadata(options.tagNamePrefix, options.tagNameSuffix)
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
          this.files[fileName].version = this.files[fileName].version + 1
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
      logger('No file to transpile')
    }

    const program = TJS.getProgramFromFiles(
      this.rootFileNames,
      { ...config.config.compilerOptions },
      this.ROOT_DIRECTORY
    )

    this.generator = TJS.buildGenerator(program, settings)
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
      getCompilationSettings: () => this.compilerOptions,
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
    try {
      const output = this.service.getEmitOutput(fileName)

      if (!output.emitSkipped) {
        if (this.verbose) {
          logger('Emit succeeded: %s', fileName)
        }
      } else {
        logger('Emit failed: %s', fileName)
        this.logErrors(fileName)
      }
    } catch (e) {
      logger('getEmitOutput failed for:\n %s \n %O \n %s', fileName, e.stack, e.message)
    }
  }

  logErrors(fileName: string) {
    const allDiagnostics = this.service
      .getCompilerOptionsDiagnostics()
      .concat(this.service.getSyntacticDiagnostics(fileName))
      .concat(this.service.getSemanticDiagnostics(fileName))

    allDiagnostics.forEach(diagnostic => {
      const message = ts.flattenDiagnosticMessageText(diagnostic.messageText, '\n')
      if (diagnostic.file) {
        const { line, character } = diagnostic.file.getLineAndCharacterOfPosition(diagnostic.start!)
        logger('  Error %j (%s, %s): %s', diagnostic.file.fileName, line + 1, character + 1, message)
      } else {
        logger('  Error: %s', message)
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
