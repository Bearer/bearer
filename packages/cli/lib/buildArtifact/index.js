const prepareConfig = require('./prepareConfig')
const attachConfig = require('./attachConfig')
const addFilesToArchive = require('./addFilesToArchive')
const generateHandler = require('./generateHandler')
const { promisify } = require('util')
const fs = require('fs-extra')
const uuidv4 = require('uuid/v4')

const readFileAsync = promisify(fs.readFile)
const writeFileAsync = promisify(fs.writeFile)
const path = require('path')

const archiver = require('archiver')

const rollup = require('rollup')
const resolve = require('rollup-plugin-node-resolve')
const commonjs = require('rollup-plugin-commonjs')
const typescript = require('rollup-plugin-typescript')
const multiEntry = require('rollup-plugin-multi-entry')
const CONFIG_FILE = 'bearer.config.json'
const HANDLER_NAME = 'index.js'
const archive = archiver('zip', {
  zlib: { level: 9 }
})

module.exports = async (
  output,
  handler,
  { path: sourceCodePath, scenarioUuid },
  emitter
) => {
  const buildFolder = path.join(path.dirname(output.path), 'build')

  const inputOptions = {
    input: path.join(buildFolder, HANDLER_NAME),
    plugins: [
      resolve({
        customResolveOptions: {
          moduleDirectory: path.join(buildFolder, 'node_modules')
        }
      }),
      commonjs(),
      typescript()
    ]
  }

  const outputOptions = {
    format: 'cjs'
  }

  if (!fs.existsSync(buildFolder)) {
    fs.mkdirSync(buildFolder)
  }

  fs.emptyDir(buildFolder)

  try {
    await fs.copy(sourceCodePath, buildFolder)
    const config = await prepareConfig(sourceCodePath, scenarioUuid)
    emitter.emit('buildArtifact:configured', { intents: config.intents })
    console.log(generateHandler(config))

    await writeFileAsync(
      path.join(buildFolder, HANDLER_NAME),
      generateHandler(config)
    )

    await writeFileAsync(
      path.join(buildFolder, CONFIG_FILE),
      JSON.stringify(config)
    )

    const bundle = await rollup.rollup(inputOptions)
    const { code, map } = await bundle.generate(outputOptions)

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

    archive.append(code, { name: HANDLER_NAME })
    archive.append(JSON.stringify(config), { name: CONFIG_FILE })
    archive.finalize()
  } catch (e) {}
  output.on('close', () => {
    emitter.emit('buildArtifact:output:close', path.resolve(output.path))
  })

  return new Promise((resolve, reject) => {
    output.on('close', () => {
      resolve(archive)
    })

    archive.on('error', reject)
  })
}
