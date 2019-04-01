import { flags } from '@oclif/command'
import getPort from 'get-port'
import * as http from 'http'
import axios from 'axios'
import * as opn from 'opn'
import * as crypto from 'crypto'

import BaseCommand from '../base-command'

// const domain = 'login.bearer.sh'

const BEARER_LOGIN_PORT = 56789
const domain = 'https://login.bearer.sh'
const callbackurl = `http://localhost:4444`

export default class Login extends BaseCommand {
  _server?: http.Server
  _verifier!: string
  _challenge!: string

  static description = 'Login to Bearer platform'

  static flags = {
    ...BaseCommand.flags,
    email: flags.string({ char: 'e' })
  }

  static args = []

  async run() {
    this._server = await this.startServer()
    this._verifier = base64URLEncode(crypto.randomBytes(32))
    this._challenge = base64URLEncode(sha256(this._verifier))
    this.ux.action.start('Logging you in')
    const scopes = 'profile+offline_access'
    const audience = `cli-${process.env.BEARER_ENV || 'production'}`
    const params = {
      audience,
      scope: scopes,
      response_type: 'code',
      client_id: this.clientId,
      code_challenge: this._challenge,
      code_challenge_method: 'S256',
      redirect_uri: callbackurl
    }
    this.debug('authoriwe params %j', params)
    const url = `${domain}/authorize?${toParams(params)}`
    opn(url)
  }

  private stopServer = () => {
    this.debug('stopping server')
    if (this._server) {
      this._server.close(() => {
        this.debug('server stopped')
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
              response.setHeader('Access-Control-Allow-Origin', this.bearerConfig.DeveloperPortalUrl)
              response.write('OK')
              response.end()
              this.getToken(data.code)
              this.stopServer()
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
      const token = await axios.post<IToken>(`${domain}/oauth/token`, {
        code,
        grant_type: 'authorization_code',
        client_id: `${this.clientId}`,
        code_verifier: `${this._verifier}`,
        redirect_uri: callbackurl
      })
      this.debug(token)
      this.ux.action.stop()
    } catch (e) {
      this.error(e)
    }
  }

  get clientId(): string {
    return 'yL4bcogYQCa7pwE2Xkc0MfDocDBBI3G3'
  }
}

function toParams(obj: Record<string, string | number>) {
  return Object.keys(obj)
    .map(key => [key, obj[key]].join('='))
    .join('&')
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

interface IToken {
  access_token: string
  refresh_token: string
  scope: string
  expires_in: number
  token_type: string
}
