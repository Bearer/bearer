const archiver = require('archiver')
const globby = require('globby')
const pathJs = require('path')
const webpack = require('webpack')

const prepareConfig = require('./prepareConfig')
const attachConfig = require('./attachConfig')
const addFilesToArchive = require('./addFilesToArchive')
const generateHandler = require('./generateHandler')

const CONFIG_FILE = 'bearer.config.json'
const HANDLER_NAME = 'index.js'
const archive = archiver('zip', {
  zlib: { level: 9 }
})

module.exports = async (output, handler, { path, scenarioUuid }, emitter) => {
  output.on('close', () => {
    emitter.emit('buildArtifact:output:close', pathJs.resolve(output.path))
  })

  output.on('end', () => {
    emitter.emit('buildArtifact:output:end')
  })

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
    .catch(console.error)

  return new Promise((resolve, reject) => {
    output.on('close', () => {
      resolve(archive)
    })

    archive.on('error', reject)
  })
}

function transpileIntents(path) {
  return new Promise((resolve, reject) => {
    globby([`${path}/*.ts`]).then(files => {
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
            // TODO: print better error messages
            console.log('[BEARER]', 'stats', stats.toJson('verbose'))
            reject(false)
          } else {
            resolve(true)
          }
        }
      )
    })
  })
}
