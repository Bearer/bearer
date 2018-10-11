import * as archiver from 'archiver'
import * as globby from 'globby'
import * as pathJs from 'path'
import * as webpack from 'webpack'

import LocationProvider from '../locationProvider'

import addFilesToArchive from './addFilesToArchive'
import attachConfig from './attachConfig'
import generateHandler from './generateHandler'
import prepareConfig from './prepareConfig'

const CONFIG_FILE = 'bearer.config.json'
const HANDLER_NAME = 'index.js'
const archive = archiver('zip', {
  zlib: { level: 9 }
})

export default (output, { scenarioUuid }, emitter, locator: LocationProvider) => {
  return new Promise(async (resolve, reject) => {
    try {
      output.on('close', () => {
        emitter.emit('buildArtifact:output:close', pathJs.resolve(output.path))
      })

      output.on('end', () => {
        emitter.emit('buildArtifact:output:end')
      })

      output.on('close', () => {
        resolve(archive)
      })

      archive.on('error', reject)

      archive.on('warning', err => {
        if (err.code === 'ENOENT') {
          emitter.emit('buildArtifact:archive:warning:ENOENT', err)
        } else {
          throw err
        }
      })

      archive.pipe(output)

      emitter.emit('buildArtifact:start', { scenarioUuid })
      // generate javascript files
      // copy package.json
      // create config
      // zip
      const distPath = locator.buildIntentsResourcePath('dist')
      await transpileIntents(locator.srcIntentsDir, distPath)

      emitter.emit('buildArtifact:intentsTranspiled')
      await prepareConfig(
        locator.authConfigPath,
        distPath,
        scenarioUuid,
        locator.scenarioRootResourcePath('node_modules'),
        locator.srcIntentsDir
      )
        .then(async config => {
          emitter.emit('buildArtifact:configured', { intents: config.intents })
          await attachConfig(archive, JSON.stringify(config, null, 2), {
            name: CONFIG_FILE
          })
          archive.append(generateHandler(config), { name: HANDLER_NAME })
          await addFilesToArchive(archive, distPath)
        })
        .then(() => {
          archive.finalize()
        })
        .catch(error => {
          emitter.emit('buildArtifact:failed', { error: error.toString() })
        })
    } catch (error) {
      emitter.emit('buildArtifact:error', error)
      reject(error)
    }
  })
}

export function transpileIntents(entriesPath: string, distPath: string): Promise<boolean | Array<any>> {
  return new Promise((resolve, reject) => {
    // Note: it works because we have client.ts present
    globby([`${entriesPath}/*.ts`])
      .then(files => {
        if (files.length) {
          const entries = files.reduceRight(
            (entriesAcc, file) => ({
              ...entriesAcc,
              [pathJs.basename(file).split('.')[0]]: file
            }),
            {}
          )
          webpack(
            {
              mode: 'production',
              // optimization: {
              //   minimize: false
              // },
              entry: entries,
              module: {
                rules: [
                  {
                    test: /\.tsx?$/,
                    loader: 'ts-loader',
                    exclude: /node_modules/,
                    options: {
                      onlyCompileBundledFiles: true,
                      compilerOptions: {
                        allowUnreachableCode: false,
                        declaration: false,
                        lib: ['es2017'],
                        noUnusedLocals: false,
                        noUnusedParameters: false,
                        allowSyntheticDefaultImports: true,
                        experimentalDecorators: true,
                        moduleResolution: 'node',
                        module: 'es6',
                        target: 'es2017'
                      }
                    }
                  }
                ]
              },
              resolve: {
                extensions: ['.tsx', '.ts', '.js']
              },
              target: 'node',
              output: {
                libraryTarget: 'commonjs2',
                filename: '[name].js',
                path: distPath
              }
              // TODO: check if it is necessary
              // context: pathJs.resolve(path)
            },
            (err, stats) => {
              if (err || stats.hasErrors()) {
                reject(stats)
              } else {
                resolve(true)
              }
            }
          )
        } else {
          reject([{ error: 'No intents to process' }])
        }
      })
      .catch(error => {
        reject([{ error }])
      })
  })
}
