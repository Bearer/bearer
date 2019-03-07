import { flags } from '@oclif/command'
import * as globby from 'globby'
import * as Listr from 'listr'
import * as path from 'path'
import * as webpack from 'webpack'

import BaseCommand from '../../base-command'
import installDependencies from '../../tasks/install-dependencies'
import { RequireIntegrationFolder } from '../../utils/decorators'
import GenerateApiDocumenation from '../generate/api-documentation'
import compilerOptions from '../../utils/intent-ts-compiler-options'

const skipInstall = 'skip-install'

export default class BuildIntents extends BaseCommand {
  static description = 'Build integration intents'
  static hidden = true
  static flags = {
    ...BaseCommand.flags,
    [skipInstall]: flags.boolean({})
  }

  static args = []

  @RequireIntegrationFolder()
  async run() {
    const { flags } = this.parse(BuildIntents)

    const tasks: Listr.ListrTask[] = [
      {
        title: 'Generate intents',
        task: async (ctx: any, _task: any) => {
          ctx.files = await this.transpile(this.locator.srcIntentsDir, this.locator.buildIntentsResourcePath('dist'))
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
      this.success('Built intents')
    } catch (e) {
      this.error(e)
    }
  }

  transpile = (entriesPath: string, distPath: string): Promise<string[]> => {
    return new Promise<string[]>(async (resolve, reject) => {
      try {
        const files = await globby([`${entriesPath}/*.ts`])
        if (!files.length) {
          return reject(new Error('No intent to transpile'))
        }

        const config: webpack.Configuration = {
          ...baseConfig,
          // optimization: {
          //   minimize: false
          // },
          entry: getEntries(files),
          output: {
            libraryTarget: 'commonjs2',
            filename: '[name].js',
            path: distPath
          }
        }

        webpack(config, (err: any, stats: webpack.Stats) => {
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
}

function getEntries(files: string[]): webpack.Entry {
  return files.reduceRight(
    (entriesAcc, file) => ({
      ...entriesAcc,
      [path.basename(file).split('.')[0]]: file
    }),
    {}
  )
}

const baseConfig: Partial<webpack.Configuration> = {
  mode: 'production',
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        loader: 'ts-loader',
        exclude: /node_modules/,
        options: {
          compilerOptions
        }
      }
    ]
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js']
  },
  target: 'node'
}
