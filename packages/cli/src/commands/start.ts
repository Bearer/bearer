import { flags } from '@oclif/command'
import * as path from 'path'
import * as chokidar from 'chokidar'
import * as http from 'http'
import * as fs from 'fs-extra'
import getPort from 'get-port'
import { spawn } from 'child_process'

import BaseCommand from '../base-command'
import GenerateApiDocumenation from './generate/api-documentation'
import GenerateSpec from './generate/spec'
import PrepareViews from './prepare/views'
import installDependencies from '../tasks/install-dependencies'
import localServerRouter from '../utils/localServerRouter'

const noOpen = 'no-open'
const noInstall = 'no-install'
// const noWatcher = 'no-watcher'
const DEFAULT_PORT = 3000

export default class Start extends BaseCommand {
  static description = 'start local development environment'
  static flags = {
    ...BaseCommand.flags,
    force: flags.boolean({ char: 'f', description: 'Start using random available port' }),
    [noOpen]: flags.boolean({}),
    [noInstall]: flags.boolean({})
    // [noWatcher]: flags.boolean({ hidden: true }) // TODO?
  }

  static args = []

  async run() {
    const { flags } = this.parse(Start)

    if (flags[noInstall]) {
      await installDependencies({ cwd: this.locator.integrationRoot })
      // TODO: install or not
    }

    const integrationHost = await this.startFunctionsServer(flags.force)

    if (this.hasViews) {
      await GenerateSpec.run(['--silent'])
      await PrepareViews.run(['--silent'])
      await GenerateApiDocumenation.run(['--silent'])
      this.startViewsServer(integrationHost)
    }
  }

  private _server!: http.Server

  stopServer = () => {
    if (this._server) {
      this._server.close()
    }
  }

  startViewsServer = async (integrationHost: string) => {
    const { flags } = this.parse(Start)
    // watch non ts files
    await this.watchNonTSFiles(this.locator.srcViewsDir, this.locator.buildViewsComponentsDir)
    const viewsVariables = {
      BEARER_INTEGRATION_ID: 'localhost',
      BEARER_INTEGRATION_HOST: integrationHost,
      BEARER_AUTHORIZATION_HOST: integrationHost
    }
    this.debug('views env variables %j', viewsVariables)
    // Build env for sub commands
    const envVariables: NodeJS.ProcessEnv = {
      ...process.env,
      ...viewsVariables
    }

    const args = [path.join(__dirname, '..', '..', 'scripts', 'startTranspiler.js')]
    const bearerTranspiler = spawn('node', args, {
      cwd: this.locator.integrationRoot,
      env: envVariables,
      stdio: ['pipe', 'pipe', 'pipe', 'ipc']
    })
    bearerTranspiler.stdout.on('data', this.childProcessStdout('transpiler'))
    bearerTranspiler.stderr.on('data', this.childProcessStderr('transpiler'))

    const tsxWatcher = chokidar.watch('**/*.tsx', {
      ignored: /(^|[\/\\])\../,
      cwd: this.locator.srcViewsDir,
      ignoreInitial: true,
      persistent: true,
      followSymlinks: false
    })
    tsxWatcher.on('add', () => bearerTranspiler.send('refresh'))
    tsxWatcher.on('unlink', () => bearerTranspiler.send('refresh'))
    tsxWatcher.on('error', error => this.warn('[transpiler] error: %j', error))

    bearerTranspiler.on('message', ({ event }) => {
      this.debug('transpiler: event received %j', event)
      if (event === 'transpiler:initialized') {
        //     // Start stencil
        const stencilServer = ['stencil', 'build', '--dev', '--watch', '--serve']
        // TODO allow yarn
        const args = [...stencilServer]
        if (flags[noOpen]) {
          args.push('--no-open')
        }

        const viewsProcess = spawn('yarn', args, {
          cwd: this.locator.buildViewsDir,
          env: envVariables
        })
        const scope = 'views'
        viewsProcess.stdout.on('data', this.childProcessStdout(scope))
        viewsProcess.stderr.on('data', this.childProcessStderr(scope))
        viewsProcess.on('close', this.childProcessClose(scope))
      }
    })
  }

  startFunctionsServer = async (force: boolean) => {
    this.debug('starting local functions server')
    const expectedPort = process.env.PORT ? parseInt(process.env.PORT, 10) : DEFAULT_PORT
    const port = await getPort({ port: expectedPort })

    // tslint:disable-next-line:no-http-string
    const host = `http://localhost:${port}`

    if (expectedPort !== DEFAULT_PORT && this.hasViews) {
      this.log('*************** Action required *****************\n')
      this.log(this.colors.red('You have specified a custom port.'))
      this.log(this.colors.yellow('You must update the views/index.html to match this setting as follow:\n'))
      this.log(`<script> bearer("CLIENT_ID", { integrationHost: "${host}" }) </script>\n`)
      this.log('*************************************************\n')
    }

    if (port !== Number(expectedPort) && !force) {
      this.error(
        `Could not start local server port ${expectedPort} is already in use.
You can specify your own port by running PORT=3322 yarn bearer start`
      )
    }

    const bearerBaseURL = `${host}/`
    process.env.bearerBaseURL = bearerBaseURL

    return new Promise<string>((resolve, reject) => {
      this._server = http.createServer(localServerRouter(this, host)).listen(port, () => {
        this.success(`Local server started. Available at this location: ${host}`)
        resolve(host)
      })
    })
  }

  childProcessStdout = (scope: string) => (data: any) => {
    this.log('[%s] %s', scope, data)
  }

  childProcessStderr = (scope: string) => (data: any) => {
    this.log(data)
  }

  childProcessClose = (scope: string) => (code: any) => {
    this.log('[%s] %n', scope, code)
  }

  watchNonTSFiles = (watchedPath: string, destPath: string): Promise<chokidar.FSWatcher> => {
    return new Promise((resolve, _reject) => {
      const callback = (error: any) => {
        if (error) {
          this.warn('error %j', error)
        }
      }
      const watcher = chokidar.watch(`${watchedPath}/**`, {
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
          this.debug('creating symlink %s %s', filePath, targetPath)
          fs.ensureSymlink(filePath, targetPath, callback)
        }
        // // Deleting symlink
        if (event === 'unlink') {
          this.debug('deleting symlink')
          fs.unlink(targetPath, err => {
            if (err) throw err
            this.debug('%s was deleted', targetPath)
          })
        }
      })
    })
  }
}
