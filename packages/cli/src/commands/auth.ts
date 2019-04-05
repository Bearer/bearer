import { flags } from '@oclif/command'
import * as http from 'http'
import * as fs from 'fs'
import axios from 'axios'
import BaseCommand from '../base-command'
// @ts-ignore
import * as open from 'open'
import getPort from 'get-port'

import { RequireIntegrationFolder } from '../utils/decorators'
import { BEARER_AUTH_PORT } from '../utils/constants'

type Event = 'success' | 'error' | 'shutdown'

export default class IntegrationsCreate extends BaseCommand {
  static description = 'create a new bearer integation'

  static flags = {
    ...BaseCommand.flags
  }

  _server?: http.Server
  _verifier!: string
  _challenge!: string
  private _listerners!: Record<Event, (() => void)[]>

  @RequireIntegrationFolder()
  async run() {
    this._listerners = {
      success: [],
      error: [],
      shutdown: []
    }
    this._server = await this.startServer()
    const config: any = JSON.parse(
      fs.readFileSync(this.locator.integrationRootResourcePath('auth.config.json'), {
        encoding: 'utf8'
      })
    )
    config.clientID = await this.askForString('Client ID', { type: 'password' })
    config.clientSecret = await this.askForString('Client secret', { type: 'password' })
    this.debug(config)
    const location = await axios
      .post(`${this.constants.IntegrationServiceHost}v2/auth/local-auth`, { config }, { maxRedirects: 0 })
      .catch(e => e.response.headers.location)

    open(`${this.constants.IntegrationServiceHost}v2/auth/${location.replace('./', '')}&clientId=NONE`)
  }

  private stopServer = () => {
    this.debug('stopping server')
    if (this._server) {
      this._server.close(() => {
        this.debug('server stopped')
        this._listerners['shutdown'].map(cb => cb())
      })
    }
  }

  private startServer = async (): Promise<http.Server> => {
    const port = await getPort({ port: BEARER_AUTH_PORT })
    return new Promise((resolve, reject) => {
      if (port !== BEARER_AUTH_PORT) {
        this.error(`bearer setup requires port ${BEARER_AUTH_PORT} to be available`)
        reject()
      }
      this.debug('starting server')
      const server = http
        .createServer((request, response) => {
          let body = ''
          request.on('data', chunk => {
            body += chunk
          })
          request.on('end', () => {
            try {
              const data: any = JSON.parse(body)
              this.debug(data)
              response.setHeader(
                'Access-Control-Allow-Origin',
                process.env.AUTH_ALLOWED_ORIGIN || this.constants.IntegrationServiceUrl
              )
              response.write('OK')
              response.end()
              this.stopServer()
            } catch (e) {
              this.debug(e)
            }
          })
        })
        .listen(BEARER_AUTH_PORT, () => {
          this.debug('server started')
          resolve(server)
        })
    })
  }
}
