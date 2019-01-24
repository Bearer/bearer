import * as path from 'path'
import * as fs from 'fs-extra'
import * as chokidar from 'chokidar'
import { spawn, execSync } from 'child_process'

import Locator from '../locationProvider'

import startLocalDevelopmentServer from './startLocalDevelopmentServer'

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
      if (event === 'add') {
        console.log('creating symlink', filePath, targetPath)
        fs.ensureSymlink(filePath, targetPath, callback)
      }
      // // Deleting symlink
      if (event === 'unlink') {
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
      const { buildViewsDir, buildViewsComponentsDir, srcViewsDir, scenarioRoot } = locator

      // Link non TS files
      const watcher = await watchNonTSFiles(srcViewsDir, buildViewsComponentsDir)
      const apiDef = locator.buildViewsResourcePath('src/openapi.json')
      if (!fs.pathExistsSync(apiDef)) {
        fs.writeJsonSync(apiDef, {})
      }
      if (!watchMode) {
        watcher.close()
      }

      if (install) {
        emitter.emit('start:prepare:installingDependencies')
        execSync(`${config.command} install`, { cwd: scenarioRoot })
      }

      return {
        rootLevel: scenarioRoot,
        buildDirectory: buildViewsDir,
        viewsDirectory: srcViewsDir
      }
    } catch (error) {
      emitter.emit('start:prepare:failed', { error })
      throw error
    }
  }
}

export const start = (emitter, config, locator: Locator) => async ({ open, install, watcher }) => {
  const { scenarioUuid } = config

  try {
    await prepare(emitter, config, locator)({
      install,
      watchMode: watcher
    })

    const { scenarioRoot, buildViewsDir } = locator
    // start local development server
    const integrationHost = await startLocalDevelopmentServer(emitter, config, locator)

    emitter.emit('start:watchers')
    if (watcher) {
      fs.watchFile(locator.authConfigPath, { persistent: true, interval: 250 }, () => {
        // TODO: ensure setup components are up to date
      })
    }

    // Start bearer transpiler phase
    const BEARER = 'bearer-transpiler'
    const options = [watcher ? null : '--no-watcher']

    // Build env for sub commands
    const envVariables: NodeJS.ProcessEnv = {
      ...process.env,
      BEARER_SCENARIO_TAG_NAME: 'localhost',
      BEARER_SCENARIO_ID: scenarioUuid,
      BEARER_INTEGRATION_HOST: integrationHost,
      BEARER_AUTHORIZATION_HOST: integrationHost
    }

    const args = [path.join(__dirname, '..', 'startTranspiler.js'), ...options].filter(el => el)
    const bearerTranspiler = spawn('node', args, {
      cwd: scenarioRoot,
      env: envVariables,
      stdio: ['pipe', 'pipe', 'pipe', 'ipc']
    })
    bearerTranspiler.stdout.on('data', childProcessStdout(emitter, BEARER))
    bearerTranspiler.stderr.on('data', childProcessStderr(emitter, BEARER))
    bearerTranspiler.on('close', childProcessClose(emitter, BEARER))

    if (watcher) {
      const tsxWatcher = chokidar.watch('**/*.tsx', {
        ignored: /(^|[\/\\])\../,
        cwd: locator.srcViewsDir,
        ignoreInitial: true,
        persistent: true,
        followSymlinks: false
      })
      tsxWatcher.on('add', () => bearerTranspiler.send('refresh'))
      tsxWatcher.on('unlink', () => bearerTranspiler.send('refresh'))
      tsxWatcher.on('error', error => emitter.emit('start:watchers:componentError', { error }))

      bearerTranspiler.on('message', ({ event }) => {
        if (event === 'transpiler:initialized') {
          // Start stencil
          const stencilServer = ['stencil', 'build', '--dev', '--watch', '--serve']
          const args = config.isYarnInstalled ? stencilServer : ['run', ...stencilServer]
          if (!open) {
            args.push('--no-open')
          }
          const stencil = spawn(config.command, args, {
            cwd: buildViewsDir,
            env: envVariables
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
