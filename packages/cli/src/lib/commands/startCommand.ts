const path = require('path')
const fs = require('fs-extra')
const copy = require('copy-template-dir')
const Case = require('case')
const chokidar = require('chokidar')
const startLocalDevelopmentServer = require('./startLocalDevelopmentServer')

const { spawn, execSync } = require('child_process')

import Locator from '../locationProvider'

function watchNonTSFiles(watchedPath, destPath): Promise<any> {
  return new Promise((resolve, _reject) => {
    function callback(error) {
      if (error) {
        console.log('error', error)
      }
    }
    const watcher = chokidar.watch(watchedPath + '/**', {
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

export function prepare(emitter, config, locator: Locator) {
  return async (
    { install = true, watchMode = true } = {
      install: true,
      watchMode: true
    }
  ) => {
    try {
      const {
        scenarioConfig: { scenarioTitle }
      } = config
      const { buildDir, srcScreenDir, buildScreenDir, scenarioRoot } = locator

      // Create hidden folder
      emitter.emit('start:prepare:buildFolder')
      if (!fs.existsSync(buildDir)) {
        fs.mkdirSync(buildDir)
        fs.mkdirSync(buildScreenDir)
      }
      fs.emptyDirSync(buildScreenDir)

      // Symlink node_modules
      emitter.emit('start:symlinkNodeModules')
      createEvenIfItExists(path.join(scenarioRoot, 'node_modules'), path.join(buildDir, 'node_modules'))

      // symlink package.json
      emitter.emit('start:symlinkPackage')

      createEvenIfItExists(path.join(scenarioRoot, 'package.json'), path.join(buildDir, 'package.json'))

      // Copy stencil.config.json
      emitter.emit('start:prepare:stencilConfig')

      const vars = {
        componentTagName: Case.kebab(scenarioTitle)
      }
      const inDir = path.join(__dirname, 'templates', 'start', '.build')
      await new Promise((resolve, reject) => {
        copy(inDir, buildDir, vars, (err, createdFiles) => {
          if (err) reject(err)
          createdFiles && createdFiles.forEach(filePath => emitter.emit('start:prepare:copyFile', filePath))
          resolve()
        })
      })

      createEvenIfItExists(path.join(buildDir, 'global'), path.join(buildScreenDir, 'global'))
      // Link non TS files
      const watcher = await watchNonTSFiles(srcScreenDir, buildScreenDir)

      if (!watchMode) {
        watcher.close()
      }

      if (install) {
        emitter.emit('start:prepare:installingDependencies')
        execSync('yarn install', { cwd: scenarioRoot })
      }

      return {
        rootLevel: scenarioRoot,
        buildDirectory: buildDir,
        screensDirectory: srcScreenDir
      }
    } catch (error) {
      emitter.emit('start:prepare:failed', { error })
      throw error
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

export const start = (emitter, config, locator: Locator) => async ({ open, install, watcher }) => {
  const {
    bearerConfig: { OrgId },
    scenarioConfig: { scenarioTitle }
  } = config

  const scenarioUuid = `${OrgId}-${scenarioTitle}`

  try {
    await prepare(emitter, config, locator)({
      install,
      watchMode: watcher
    })

    const { scenarioRoot, buildDir } = locator
    /* start local development server */
    const integrationHost = await startLocalDevelopmentServer(scenarioRoot, scenarioUuid, emitter, config)

    ensureSetupAndConfigComponents(buildDir)

    emitter.emit('start:watchers')
    if (watcher) {
      fs.watchFile(path.join(scenarioRoot, 'auth.config.json'), { persistent: true, interval: 250 }, () =>
        ensureSetupAndConfigComponents(buildDir)
      )
    }

    /* Start bearer transpiler phase */
    const BEARER = 'bearer-transpiler'
    const bearerTranspiler = spawn(
      'node',
      [path.join(__dirname, '..', 'startTranspiler.js'), watcher ? null : '--no-watcher'].filter(el => el),
      {
        cwd: scenarioRoot,
        env: {
          ...process.env,
          BEARER_SCENARIO_ID: scenarioUuid,
          BEARER_INTEGRATION_HOST: integrationHost
        },
        stdio: ['pipe', 'pipe', 'pipe', 'ipc']
      }
    )
    bearerTranspiler.stdout.on('data', childProcessStdout(emitter, BEARER))
    bearerTranspiler.stderr.on('data', childProcessStderr(emitter, BEARER))
    bearerTranspiler.on('close', childProcessClose(emitter, BEARER))

    if (watcher) {
      bearerTranspiler.on('message', ({ event }) => {
        if (event === 'transpiler:initialized') {
          /* Start stencil */
          const args = ['start']
          if (!open) {
            args.push('--no-open')
          }
          const stencil = spawn('yarn', args, {
            cwd: buildDir,
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

export function useWith(program, emitter, config, locator: Locator) {
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
    .action(start(emitter, config, locator))
}

/**
 * Logger helpers
 */

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

function createEvenIfItExists(target, sourcePath): void {
  try {
    fs.symlinkSync(target, sourcePath)
  } catch (e) {
    if (e.code !== 'EEXIST') {
      throw e
    }
  }
}
