import * as http from 'http'
import * as fs from 'fs'
import axios from 'axios'

import { TConfig, Authentications, configs } from '@bearer/types/lib/authentications'

import BaseCommand from '../../base-command'
// @ts-ignore
import * as open from 'open'
import getPort from 'get-port'

import { RequireIntegrationFolder } from '../../utils/decorators'
import { BEARER_AUTH_PORT } from '../../utils/constants'

type Event = 'success' | 'error' | 'shutdown'

export default class SetupAuth extends BaseCommand {
  static description = 'setup local credentials'

  static flags = {
    ...BaseCommand.flags
  }

  _server?: http.Server
  _verifier!: string
  _challenge!: string
  private _listerners!: Record<Event, ((data: any) => void)[]>

  @RequireIntegrationFolder()
  async run() {
    this._listerners = {
      success: [],
      error: [],
      shutdown: []
    }
    const config: TConfig = JSON.parse(
      fs.readFileSync(this.locator.integrationRootResourcePath('auth.config.json'), {
        encoding: 'utf8'
      })
    )
    const { authType } = config
    switch (authType) {
      case Authentications.OAuth2: {
        const { token } = await this.fetchAuthTken(config as configs.TOAuth2Config)
        this.log('Your token: %s', token)
        // TODO: save
        break
      }

      case Authentications.Basic: {
        const username = await this.askForString('Username')
        const password = await this.askForPassword('Password')
        this.debug(username, password)
        // TODO: save
        this.log('Your credentials: username => %s  password => ', username, password.replace(/./g, '*'))
        break
      }
      case Authentications.ApiKey: {
        const apiKey = await this.askForPassword('API Key')
        this.debug(apiKey)
        // TODO: save
        this.log('Your api Key: %s', apiKey)
        break
      }
      case Authentications.OAuth1: {
        break
      }
      case Authentications.Custom:
      case Authentications.NoAuth: {
        return this.warn(
          `The current authentication type of the integraion is not supported by this command: ${authType}`
        )
      }
      default:
        // unsure we handled all authentications
        // http://ideasintosoftware.com/exhaustive-switch-in-typescript/
        this.error(`The current authentication type of the integraion is not supported by this command: ${authType}`)
        throw new UnreachableCaseError(authType)
    }
  }

  fetchAuthTken = async (config: configs.TOAuth2Config): Promise<{ token: string }> => {
    return new Promise(async (resolve, reject) => {
      this._server = await this.startServer()
      const { BEARER_AUTH_CLIENT_ID, BEARER_AUTH_CLIENT_SECRET } = process.env
      const clientID = BEARER_AUTH_CLIENT_ID || (await this.askForString('Client ID', { type: 'password' }))
      const clientSecret = BEARER_AUTH_CLIENT_SECRET || (await this.askForString('Client secret', { type: 'password' }))
      const newConfig = { ...config, clientID, clientSecret }
      const location = await axios
        .post(`${this.constants.IntegrationServiceHost}v2/auth/local-auth`, { config: newConfig }, { maxRedirects: 0 })
        .catch(e => e.response.headers.location)
      this.debug(config)

      this.on<{ token: string }>('success', data => {
        resolve(data)
      })

      this.on('error', data => {
        this.debug(data)
        reject('Error while receiving token')
      })
      const url = `${this.constants.IntegrationServiceHost}v2/auth/${location.replace('./', '')}&clientId=NONE`
      const a = await open(url)
      a.on('close', (code: any, signal: any) => {
        if (code !== 0) {
          this.stopServer()
          this.warn(
            this.colors.yellow(`Unable to open a browser. If you want to retrieve a token please follow these steps\n`)
          )
          this.log(this.colors.bold('1/ access the url below  and follow the login process:\n\n'), url)
          this.log()
          this.log(this.colors.bold(`2/ when you access the success page copy the token and paste it here`))
          this.askForString('Token').then(token => {
            resolve({ token })
          })
        }
      })
    })
  }

  private stopServer = () => {
    this.debug('stopping server')
    if (this._server) {
      this._server.close(() => {
        this.debug('server stopped')
        this._listerners['shutdown'].map(cb => cb({}))
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
              const data: { token: string } = JSON.parse(body)
              this.debug(data)
              response.setHeader(
                'Access-Control-Allow-Origin',
                process.env.AUTH_ALLOWED_ORIGIN || this.constants.IntegrationServiceUrl
              )
              response.setHeader('Connection', 'close')
              response.write('OK')
              response.end()
              this._listerners['success'].map(cb => cb(data))
            } catch (e) {
              this._listerners['error'].map(cb => cb(e))
            }
            this.stopServer()
          })
        })
        .listen(BEARER_AUTH_PORT, () => {
          this.debug('server started')
          resolve(server)
        })
    })
  }

  on = <T>(event: Event, callback: (data: T) => void) => {
    this._listerners[event] = [...this._listerners[event], callback as any]
  }
}

class UnreachableCaseError extends Error {
  constructor(val: never) {
    super(`Unreachable case: ${val}`)
  }
}
