import { flags } from '@oclif/command'
import * as http from 'http'
import getPort from 'get-port'

import BaseCommand from '../base-command'
import GenerateApiDocumenation from './generate/api-documentation'
import GenerateSpec from './generate/spec'
import PrepareViews from './prepare/views'

const noOpen = 'no-open'
const noInstall = 'no-install'
const noWatcher = 'no-watcher'
const DEFAULT_PORT = 3000

export default class Start extends BaseCommand {
  static description = 'start local development environment'
  static server: http.Server
  static flags = {
    ...BaseCommand.flags,
    force: flags.boolean({ char: 'f', description: 'Start using random available port' }),
    [noOpen]: flags.boolean({}),
    [noInstall]: flags.boolean({}),
    [noWatcher]: flags.boolean({ hidden: true })
  }

  static args = []

  async run() {
    const { flags } = this.parse(Start)
    const cmdArgs = []
    if (flags[noOpen]) {
      cmdArgs.push(`--${noOpen}`)
    }
    if (flags[noInstall]) {
      cmdArgs.push(`--${noInstall}`)
    }
    if (flags[noWatcher]) {
      cmdArgs.push(`--${noWatcher}`)
    }

    if (flags.force) {
      cmdArgs.push(`--force`)
    }

    if (this.hasViews) {
      await GenerateSpec.run(['--silent'])
      await PrepareViews.run(['--silent'])
      await GenerateApiDocumenation.run(['--silent'])
    } else {
      cmdArgs.push(`--no-views`)
    }

    await this.startFunctionsServer(flags.force)

    // this.runLegacy(['start', ...cmdArgs])
  }

  startFunctionsServer = async (force: boolean) => {
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

    return new Promise<http.Server>((resolve, reject) => {
      const _server = http
        .createServer((request, response) => {
          this.debug('Incoming request')
          let body = ''
          request.on('data', chunk => {
            body += chunk
          })
          request.on('end', async () => {
            try {
              this.debug('method: %s body: %s', request.method, body || '{}')
              response.setHeader('Connection', 'close')
              // handle response
            } catch (e) {}
            response.end()
          })
        })
        .listen(port, () => {
          this.success('Local server started')
          resolve(_server)
        })
    })
  }
}
