import archiver from 'archiver'
import globby from 'globby'
import pathJs from 'path'
import webpack from 'webpack'

import prepareConfig from './prepareConfig'
import attachConfig from './attachConfig'
import addFilesToArchive from './addFilesToArchive'
import generateHandler from './generateHandler'

const CONFIG_FILE = 'bearer.config.json'
const HANDLER_NAME = 'index.js'
const archive = archiver('zip', {
  zlib: { level: 9 }
})

export default (output, { path, scenarioUuid }, emitter) => {
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

      await transpileIntents(path)

      await prepareConfig(path, scenarioUuid)
        .then(async config => {
          emitter.emit('buildArtifact:configured', { intents: config.intents })
          await attachConfig(archive, JSON.stringify(config, null, 2), {
            name: CONFIG_FILE
          })
          archive.append(generateHandler(config), { name: HANDLER_NAME })
          await addFilesToArchive(archive, path)
        })
        .then(() => {
          archive.finalize()
        })
        .catch(error => {
          emitter.emit('buildArtifact:failed', { error })
        })
    } catch (error) {
      emitter.emit('buildArtifact:error', error)
    }
  })
}

function transpileIntents(path) {
  return new Promise((resolve, reject) => {
    // Note: it works because we have client.ts present
    globby([`${path}/*.ts`])
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
              entry: entries,
              module: {
                rules: [
                  {
                    test: /\.tsx?$/,
                    loader: 'ts-loader',
                    exclude: /node_modules/,
                    options: {
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
                path: pathJs.join(pathJs.resolve(path), 'dist')
              },
              context: pathJs.resolve(path)
            },
            (err, stats) => {
              if (err || stats.hasErrors()) {
                reject(stats.toJson('verbose').errors)
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
