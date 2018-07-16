const path = require('path')
const fs = require('fs-extra')
const copy = require('copy-template-dir')
const Case = require('case')
const chokidar = require('chokidar')
const startLocalDevelopmentServer = require('./startLocalDevelopmentServer')

const { spawn, execSync } = require('child_process')

function createEvenIfItExists(target, sourcePath) {
  try {
    fs.symlinkSync(target, sourcePath)
  } catch (e) {
    if (!e.code === 'EEXIST') {
      throw e
    }
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

function watchNonTSFiles(watchedPath, destPath) {
  return new Promise((resolve, reject) => {
    function callback(error) {
      if (error) {
        console.log('error', error)
      }
    }
    watcher = chokidar.watch(watchedPath + '/**', {
      ignored: /\.tsx?$/,
      persistent: true,
      followSymlinks: false
    })

    watcher.on('ready', () => {
      resolve(watcher)
    })

    watcher.on('all', (event, filePath) => {
      const relativePath = filePath.replace(watchedPath, '')
      const targetPath = path.join(destPath, relativePath)
      // Creating symlink
      if (event == 'add') {
        console.log('creating symlink', filePath, targetPath)
        fs.ensureSymlink(filePath, targetPath, callback)
      }
      // // Deleting symlink
      if (event == 'unlink') {
        console.log('deleting symlink')
        fs.unlink(targetPath, err => {
          if (err) throw err
          console.log(targetPath + ' was deleted')
        })
      }
    })
  })
}

function prepare(emitter, config) {
  return async (
    { install = true, watchMode = true } = { install: true, watchMode: true }
  ) => {
    try {
      const {
        rootPathRc,
        scenarioConfig: { scenarioTitle }
      } = config

      const rootLevel = path.dirname(rootPathRc)
      const screensDirectory = path.join(rootLevel, 'screens')
      const buildDirectory = path.join(screensDirectory, '.build')
      const buildSrcDirectory = path.join(buildDirectory, 'src')

      // Create hidden folder
      emitter.emit('start:prepare:buildFolder')
      if (!fs.existsSync(buildDirectory)) {
        fs.mkdirSync(buildDirectory)
        fs.mkdirSync(buildSrcDirectory)
      }
      fs.emptyDirSync(buildSrcDirectory)

      // Symlink node_modules
      emitter.emit('start:symlinkNodeModules')
      createEvenIfItExists(
        path.join(screensDirectory, 'node_modules'),
        path.join(buildDirectory, 'node_modules')
      )

      // symlink package.json
      emitter.emit('start:symlinkPackage')

      createEvenIfItExists(
        path.join(screensDirectory, 'package.json'),
        path.join(buildDirectory, 'package.json')
      )

      // Copy stencil.config.json
      emitter.emit('start:prepare:stencilConfig')

      const vars = {
        componentTagName: Case.kebab(scenarioTitle)
      }
      const inDir = path.join(__dirname, 'templates/start/.build')
      const outDir = buildDirectory

      await new Promise((resolve, reject) => {
        copy(inDir, outDir, vars, (err, createdFiles) => {
          if (err) reject(err)
          createdFiles.forEach(filePath =>
            emitter.emit('start:prepare:copyFile', filePath)
          )
          resolve()
        })
      })

      createEvenIfItExists(
        path.join(buildDirectory, 'global'),
        path.join(buildSrcDirectory, 'global')
      )

      // Link non TS files
      const watcher = await watchNonTSFiles(
        path.join(screensDirectory, 'src'),
        path.join(buildDirectory, 'src')
      )

      if (!watchMode) {
        watcher.close()
      }

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

const ensureSetupAndConfigComponents = rootLevel => {
  spawn('bearer', ['g', '--config'], {
    cwd: rootLevel
  })
  spawn('bearer', ['g', '--setup'], {
    cwd: rootLevel
  })
}

const start = (emitter, config) => async ({ open, install, watcher }) => {
  const {
    bearerConfig: { OrgId },
    scenarioConfig: { scenarioTitle }
  } = config

  const scenarioUuid = `${OrgId}-${scenarioTitle}`

  try {
    const { buildDirectory, rootLevel, screensDirectory } = await prepare(
      emitter,
      config
    )({
      install,
      watchMode: watcher
    })

    ensureSetupAndConfigComponents(rootLevel)

    emitter.emit('start:watchers')
    if (watcher) {
      fs.watchFile(
        path.join(rootLevel, 'intents', 'auth.config.json'),
        { persistent: true, interval: 250 },
        () => ensureSetupAndConfigComponents(rootLevel)
      )
    }

    /* Start bearer transpiler phase */
    const BEARER = 'bearer-transpiler'
    const bearerTranspiler = spawn(
      'node',
      [
        path.join(__dirname, '..', 'startTranspiler.js'),
        watcher ? null : '--no-watcher'
      ].filter(el => el),
      {
        cwd: screensDirectory,
        env: {
          ...process.env,
          BEARER_SCENARIO_ID: scenarioUuid
        },
        stdio: ['pipe', 'pipe', 'pipe', 'ipc']
      }
    )
    bearerTranspiler.stdout.on('data', childProcessStdout(emitter, BEARER))
    bearerTranspiler.stderr.on('data', childProcessStderr(emitter, BEARER))
    bearerTranspiler.on('close', childProcessClose(emitter, BEARER))

    if (watcher) {
      /* start local development server */
      const { host, port } = await startLocalDevelopmentServer(
        rootLevel,
        scenarioUuid,
        emitter,
        config
      )
      const integrationHost = `http://${host}:${port}`

      bearerTranspiler.on('message', ({ event }) => {
        if (event === 'transpiler:initialized') {
          /* Start stencil */
          const args = ['start']
          if (!open) {
            args.push('--no-open')
          }
          const stencil = spawn('yarn', args, {
            cwd: buildDirectory,
            env: {
              ...process.env,
              BEARER_SCENARIO_ID: scenarioUuid,
              BEARER_INTEGRATION_HOST: integrationHost
            }
          })

          const STENCIL = 'stencil'

          stencil.stdout.on('data', childProcessStdout(emitter, STENCIL))
          stencil.stderr.on('data', childProcessStderr(emitter, STENCIL))
          stencil.on('close', childProcessClose(emitter, STENCIL))
        }
      })
    }
  } catch (e) {
    emitter.emit('start:failed', { error: e })
  }
}

module.exports = {
  start,
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
      .option('--no-watcher', 'Run transpiler only once')
      .action(start(emitter, config))
  }
}
