import getPort from 'get-port'
import * as http from 'http'
import axios from 'axios'
// @ts-ignore
import * as opn from 'open'
import * as crypto from 'crypto'

import BaseCommand from '../base-command'
import { TAccessToken } from '../types'
import { toParams } from '../utils/helpers'
import { LOGIN_CLIENT_ID, BEARER_ENV } from '../utils/constants'

const BEARER_LOGIN_PORT = 56789
type Event = 'success' | 'error' | 'shutdown'

export default class Login extends BaseCommand {
  _server?: http.Server
  _verifier!: string
  _challenge!: string
  private _listerners!: Record<Event, (() => void)[]>

  static description = 'Login to Bearer platform'

  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  async run() {
    this._listerners = {
      success: [],
      error: [],
      shutdown: []
    }
    this._server = await this.startServer()
    this._verifier = base64URLEncode(crypto.randomBytes(32))
    this._challenge = base64URLEncode(sha256(this._verifier))
    this.ux.action.start('Logging you in')

    const scopes = 'offline_access email openid'
    const audience = `cli-${BEARER_ENV}`
    const params = {
      audience,
      scope: scopes,
      response_type: 'code',
      client_id: LOGIN_CLIENT_ID,
      code_challenge: this._challenge,
      code_challenge_method: 'S256',
      redirect_uri: this.callbackUrl
    }
    this.debug('authoriwe params %j', params)
    const url = `${this.constants.LoginDomain}/authorize?${toParams(params)}`
    opn(url)

    await Promise.all([
      new Promise((resolve, reject) => {
        this.on('success', resolve)
        this.on('error', reject)
      }),
      new Promise((resolve, reject) => {
        this.on('shutdown', resolve)
        this.on('error', reject)
      })
    ])
  }

  on = (event: Event, callback: () => void) => {
    this._listerners[event] = [...this._listerners[event], callback]
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
    const port = await getPort({ port: BEARER_LOGIN_PORT })
    return new Promise((resolve, reject) => {
      if (port !== BEARER_LOGIN_PORT) {
        this.error(`bearer login requires port ${BEARER_LOGIN_PORT} to be available`)
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
              const data: { code: string } = JSON.parse(body)
              this.debug(data)
              response.setHeader(
                'Access-Control-Allow-Origin',
                process.env.LOGIN_ALLOWED_ORIGIN || this.constants.DeveloperPortalUrl
              )
              response.write('OK')
              response.end()
              this.stopServer()
              this.getToken(data.code)
            } catch (e) {
              this.debug(e)
            }
          })
        })
        .listen(BEARER_LOGIN_PORT, () => resolve(server))
    })
  }

  getToken = async (code: string) => {
    try {
      const { data: token } = await axios.post<TAccessToken>(`${this.constants.LoginDomain}/oauth/token`, {
        code,
        grant_type: 'authorization_code',
        client_id: `${LOGIN_CLIENT_ID}`,
        code_verifier: `${this._verifier}`,
        redirect_uri: this.callbackUrl
      })

      this.debug(token)
      await this.bearerConfig.storeToken(token)
      this.ux.action.stop()
      this.success('Successfully logged in!! 🐻')
      this._listerners['success'].map(cb => cb())
    } catch (e) {
      this.error(e)
    }
  }

  get callbackUrl(): string {
    return process.env.BEARER_LOGIN_CALLBACK_URL || `${this.constants.DeveloperPortalUrl}cli-login-callback`
  }
}

function base64URLEncode(str: Buffer) {
  return str
    .toString('base64')
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '')
}

function sha256(str: string) {
  return crypto
    .createHash('sha256')
    .update(str)
    .digest()
}
