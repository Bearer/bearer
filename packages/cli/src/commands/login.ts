import * as http from 'http'
import * as express from 'express'
import * as opn from 'open'
import * as crypto from 'crypto'
import getPort from 'get-port'
import axios from 'axios'
// @ts-ignore
import cliUx from 'cli-ux'

import BaseCommand from '../base-command'

import { TAccessToken } from '../types'
import { toParams } from '../utils/helpers'
import { LOGIN_CLIENT_ID, BEARER_ENV, BEARER_LOGIN_PORT, SUCCESS_LOGIN_PAGE } from '../utils/constants'
import { askForString } from '../utils/prompts'

type Event = 'success' | 'error' | 'shutdown'

export default class Login extends BaseCommand {
  static description = 'login using Bearer credentials'

  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  _server?: http.Server
  _verifier!: string
  _challenge!: string
  private _listerners!: Record<Event, (() => void)[]>

  async run() {
    this._listerners = {
      success: [],
      error: [],
      shutdown: []
    }
    this._server = await this.startServer()
    this._verifier = base64URLEncode(crypto.randomBytes(32))
    this._challenge = base64URLEncode(sha256(this._verifier))
    cliUx.action.start('Logging you in')

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
    const spawned = await opn(url)
    await Promise.race([
      new Promise((resolve, reject) => {
        spawned.on('close', async (code: any, signal: any) => {
          if (code !== 0) {
            this.stopServer()
            this.warn(
              this.colors.yellow(
                `Unable to open a browser. If you want to retrieve a token please follow these steps\n`
              )
            )
            this.log(this.colors.bold('1/ access the url below  and follow the login process:\n\n'), url)
            this.log()
            this.log(this.colors.bold(`2/ when you access the success page copy the token and paste it here`))
            const token = await askForString('Token')
            await this.getToken(token)
          }
        })
      }),
      Promise.all([
        new Promise((resolve, reject) => {
          this.on('success', resolve)
          this.on('error', reject)
        }),
        new Promise((resolve, reject) => {
          this.on('shutdown', resolve)
          this.on('error', reject)
        })
      ])
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
        this._listerners.shutdown.map(cb => cb())
      })
    }
    cliUx.action.stop()
  }

  private startServer = async (): Promise<http.Server> => {
    const port = await getPort({ port: BEARER_LOGIN_PORT })
    return new Promise((resolve, reject) => {
      if (port !== BEARER_LOGIN_PORT) {
        this.error(`bearer login requires port ${BEARER_LOGIN_PORT} to be available`)
        reject()
      }
      this.debug('starting server')
      const app = express()
      app.use((_req: express.Request, res: express.Response, next: express.NextFunction) => {
        res.setHeader('Connection', 'close')
        next()
      })
      app.get('/login/callback', (req: express.Request, res: express.Response) => {
        const code: string = req.query.code || ''
        try {
          res.send(SUCCESS_LOGIN_PAGE)
          res.end()
          this.getToken(code)
        } catch (e) {
          this.debug(e)
          this.error('Error while fetching token')
        }
        this.stopServer()
      })

      const server = app.listen(BEARER_LOGIN_PORT, () => {
        this.debug('server started')
        resolve(server)
      })
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

      this.debug('token: %j', token)
      await this.bearerConfig.storeToken(token)
      this.success('Successfully logged in!! 🐻')
      this._listerners.success.map(cb => cb())
    } catch (e) {
      this.error(e)
    }
  }

  get callbackUrl(): string {
    return process.env.BEARER_LOGIN_CALLBACK_URL || `${this.constants.DeveloperPortalUrl}cli/login-callback`
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
