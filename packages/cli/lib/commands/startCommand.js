const path = require('path')
const fs = require('fs')
const copy = require('copy-template-dir')
const Case = require('case')
const { spawn, execSync } = require('child_process')

function prepare(emitter, config) {
  return async ({ install = true } = { install: true }) => {
    try {
      const {
        rootPathRc,
        scenarioConfig: { scenarioTitle }
      } = config
      const rootLevel = path.dirname(rootPathRc)
      const screensDirectory = path.join(rootLevel, 'screens')
      const buildDirectory = path.join(screensDirectory, '.build')

      // Create hidden folder
      emitter.emit('start:prepare:buildFolder')
      if (!fs.existsSync(buildDirectory)) {
        fs.mkdirSync(buildDirectory)
        fs.mkdirSync(path.join(buildDirectory, 'src'))
      }

      // Symlink node_modules
      emitter.emit('start:symlinkNodeModules')
      const nodeModuleLink = path.join(buildDirectory, 'node_modules')
      createEvenIfItExists(
        path.join(screensDirectory, 'node_modules'),
        nodeModuleLink
      )

      // symlink package.json
      emitter.emit('start:symlinkPackage')

      const packageLink = path.join(buildDirectory, 'package.json')
      createEvenIfItExists(
        path.join(screensDirectory, 'package.json'),
        packageLink
      )

      // Copy stencil.config.json
      emitter.emit('start:prepare:stencilConfig')

      const vars = {
        componentTagName: Case.kebab(scenarioTitle)
      }
      const inDir = path.join(__dirname, 'templates/start/.build')
      const outDir = buildDirectory

      copy(inDir, outDir, vars, (err, createdFiles) => {
        if (err) throw err
        createdFiles.forEach(filePath =>
          emitter.emit('start:prepare:copyFile', filePath)
        )
      })

      if (install) {
        emitter.emit('start:prepare:installingDependencies')
        execSync('yarn install', { cwd: screensDirectory })
      }

      return {
        rootLevel,
        buildDirectory,
        screensDirectory
      }
    } catch (error) {
      emitter.emit('start:prepare:failed', { error })
      return {}
    }
  }
}

const start = (emitter, config) => async ({ open, install }) => {
  try {
    const { buildDirectory, rootLevel } = await prepare(emitter, config)({
      install
    })
    emitter.emit('start:watchers')

    /* Start bearer transpiler phase */
    const bearerTranspiler = spawn(
      'node',
      [path.join(__dirname, '..', 'startTranspiler.js')],
      {
        cwd: rootLevel
      }
    )

    const BEARER = 'bearer-transpiler'
    bearerTranspiler.stdout.on('data', childProcessStdout(emitter, BEARER))
    bearerTranspiler.stderr.on('data', childProcessStderr(emitter, BEARER))
    bearerTranspiler.on('close', childProcessClose(emitter, BEARER))

    /* Start stencil */
    const args = ['start']
    if (!open) {
      args.push('--no-open')
    }
    const stencil = spawn('yarn', args, {
      cwd: buildDirectory
    })

    const STENCIL = 'stencil'

    stencil.stdout.on('data', childProcessStdout(emitter, STENCIL))
    stencil.stderr.on('data', childProcessStderr(emitter, STENCIL))
    stencil.on('close', childProcessClose(emitter, STENCIL))
  } catch (e) {
    emitter.emit('start:failed', { error: e })
  }
}

function childProcessStdout(emitter, name) {
  return data => {
    emitter.emit('start:watchers:stdout', { name, data })
  }
}

function childProcessStderr(emitter, name) {
  return data => {
    emitter.emit('start:watchers:stderr', { name, data })
  }
}

function childProcessClose(emitter, name) {
  return code => {
    emitter.emit('start:watchers:close', { name, code })
  }
}

function createEvenIfItExists(target, path) {
  try {
    fs.symlinkSync(target, path)
  } catch (e) {
    if (!e.code === 'EEXIST') {
      throw e
    }
  }
}

module.exports = {
  prepare,
  useWith: (program, emitter, config) => {
    program
      .command('start')
      .description(
        `Start local development server.
    $ bearer start
`
      )
      .option('--no-open', 'Do not open web browser')
      .option('--no-install', 'Do not run yarn|npm install')
      .action(start(emitter, config))
  }
}
