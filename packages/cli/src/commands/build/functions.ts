import { flags } from '@oclif/command'
import * as globby from 'globby'
import * as fs from 'fs-extra'
import * as ts from 'typescript'
import * as Listr from 'listr'
import * as path from 'path'
import * as webpack from 'webpack'

import BaseCommand from '../../base-command'
import installDependencies from '../../tasks/install-dependencies'
import { RequireIntegrationFolder } from '../../utils/decorators'
import GenerateApiDocumenation from '../generate/api-documentation'
import compilerOptions from '../../utils/function-ts-compiler-options'
import prepareConfig, { HANDLER_NAME_WITH_EXT } from '../../utils/prepare-config'

const skipInstall = 'skip-install'

export default class BuildFunctions extends BaseCommand {
  static description = 'Build integration functions'
  static aliases = ['b:f']
  static hidden = true
  static flags = {
    ...BaseCommand.flags,
    [skipInstall]: flags.boolean({})
  }

  static args = []

  @RequireIntegrationFolder()
  async run() {
    const { flags } = this.parse(BuildFunctions)

    const tasks: Listr.ListrTask[] = [
      {
        title: 'Generate functions',
        task: async (ctx: any, _task: any) => {
          ctx.files = await this.transpile(
            this.locator.srcFunctionsDir,
            this.locator.buildFunctionsResourcePath('dist')
          )
        }
      },
      {
        title: 'Generate openapi.json',
        task: async (_ctx: any, _task: any) => {
          await GenerateApiDocumenation.run([])
        }
      }
    ]
    if (!flags[skipInstall]) {
      tasks.unshift(installDependencies({ cwd: this.locator.integrationRoot }))
    }

    try {
      const ctx = await new Listr(tasks).run()
      this.debug('Transpiled :\n', ctx.files.join('\n  * '))
      this.success('Built functions')
    } catch (e) {
      this.error(e)
    }
  }

  transpile = (entriesPath: string, distPath: string): Promise<string[]> => {
    return new Promise<string[]>(async (resolve, reject) => {
      try {
        // generate bundle
        const files = await globby([`${entriesPath}/*.ts`])

        if (!files.length) {
          return reject(new Error('No func to transpile'))
        }

        // transpile ts files
        files.forEach(file => {
          const content = fs.readFileSync(file, { encoding: 'utf8' })
          const { outputText } = ts.transpileModule(content, {
            compilerOptions: { ...compilerOptions, module: ts.ModuleKind.ES2015 }
          })
          const newFile = path.join(this.locator.buildFunctionsDir, file.replace(entriesPath, '').replace(/ts$/, 'js'))
          fs.writeFileSync(newFile, outputText, {
            encoding: 'utf8'
          })
        })

        const functions = await this.retrieveFunctions()
        const indexHandler = this.locator.buildFunctionsResourcePath(HANDLER_NAME_WITH_EXT)
        // generate handler

        fs.writeFileSync(indexHandler, buildLambdaIndex(functions), {
          encoding: 'utf8'
        })

        const config: webpack.Configuration = {
          ...baseConfig,
          optimization: {
            minimize: false
          },
          entry: indexHandler,
          output: {
            libraryTarget: 'commonjs2',
            filename: HANDLER_NAME_WITH_EXT,
            path: distPath
          },
          // the sdk is already provided within the lamnda through a layer
          externals: /aws\-sdk/
        }

        webpack(config, (err: any, stats: webpack.Stats) => {
          this.debug(stats.toJson('verbose'))
          if (err || stats.hasErrors()) {
            reject(
              stats.toString({
                builtAt: false,
                entrypoints: false,
                assets: false,
                version: false,
                timings: false,
                hash: false,
                modules: false,
                chunks: false, // Makes the build much quieter
                colors: true // Shows colors in the console
              })
            )
          } else {
            resolve(files)
          }
        })
      } catch (e) {
        reject(e)
      }
    })
  }

  async retrieveFunctions(): Promise<string[]> {
    try {
      const config = await prepareConfig(
        this.locator.authConfigPath,
        this.bearerConfig.bearerUid,
        this.locator.srcFunctionsDir
      )
      return config.functions
    } catch (e) {
      throw e
    }
  }
}

function buildLambdaIndex(functions: string[]): string {
  return functions.reduce(
    (out, func, index) => {
      const funcConstName = `func${index}`
      return `${out}\n
const ${funcConstName} = require("./${func}").default;
module.exports['${func}'] = ${funcConstName}.init();
`
    },
    `const bearerOverride = require('@bearer/x-ray').bearerOverride;
bearerOverride()
// functions
`
  )
}

const baseConfig: Partial<webpack.Configuration> = {
  mode: 'production',
  resolve: {
    extensions: ['.js', '.json']
  },
  target: 'node'
}
