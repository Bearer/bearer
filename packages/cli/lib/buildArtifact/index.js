const prepareConfig = require('./prepareConfig')
const generateHandler = require('./generateHandler')
const { promisify } = require('util')
const fs = require('fs')
const ZipPlugin = require('zip-webpack-plugin')
const webpack = require('webpack')

const readFileAsync = promisify(fs.readFile)
const writeFileAsync = promisify(fs.writeFile)
const pathJs = require('path')

const CONFIG_FILE = 'bearer.config.json'
const HANDLER_NAME = 'index.js'

class AddConfigWebpackPlugin {
  constructor(options) {
    this.options = options
  }
  apply(compiler) {
    const { config } = this.options
    compiler.hooks.thisCompilation.tap('emit', compilation => {
      compilation.assets[CONFIG_FILE] = {
        source() {
          return config
        },
        size() {
          return config.length
        }
      }
    })
  }
}

module.exports = async (output, handler, { path, scenarioUuid }, emitter) => {
  const artifactDirectory = pathJs.dirname(pathJs.resolve(output.path))
  const buildDirectory = pathJs.join(artifactDirectory, 'build')
  const handlerBuildPath = pathJs.join(buildDirectory, HANDLER_NAME)
  let config

  if (!fs.existsSync(buildDirectory)) {
    fs.mkdirSync(buildDirectory)
  }

  try {
    config = await prepareConfig(path, scenarioUuid)
    emitter.emit('buildArtifact:configured', { intents: config.intents })
    await writeFileAsync(handlerBuildPath, generateHandler(config))
  } catch (e) {
    console.log(e)
  }

  return new Promise((resolve, reject) => {
    webpack(
      {
        mode: 'production',
        target: 'node',
        entry: handlerBuildPath,
        output: {
          libraryTarget: 'commonjs',
          path: artifactDirectory,
          filename: 'index.js'
        },
        resolve: {
          extensions: ['.ts', '.tsx', '.js']
        },
        module: {
          rules: [
            {
              test: /\.tsx?$/,
              loader: 'awesome-typescript-loader',
              exclude: [
                pathJs.join(pathJs.resolve(path), 'node_modules'),
                /node_modules/
              ]
            }
          ]
        },
        stats: {
          colors: true,
          modules: true,
          reasons: true,
          errorDetails: true
        },
        plugins: [
          new AddConfigWebpackPlugin({
            config: JSON.stringify(config, null, 2)
          }),
          new ZipPlugin({
            path: artifactDirectory,
            filename: `${scenarioUuid}.zip`
          })
        ]
      },
      (err, stats) => {
        if (err || stats.hasErrors()) {
          console.log(err)
          console.log(stats.compilation.errors)
          reject(err, stats)
          // Handle errors here
        } else {
          resolve(stats)
        }
      }
    )
  })
}
