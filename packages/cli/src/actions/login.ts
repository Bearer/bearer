import * as http from 'http'
import * as express from 'express'
import * as crypto from 'crypto'
import getPort from 'get-port'
import axios from 'axios'
import baseCommand from '../base-command'
import * as inquirer from 'inquirer'

// TODO: use esModuleInterop config
// @ts-ignore
import cliUx from 'cli-ux'
// @ts-ignore
import * as opn from 'open'

import BaseAction, { createExport } from './base'

import { TAccessToken } from '../types'
import { toParams } from '../utils/helpers'
import { LOGIN_CLIENT_ID, BEARER_ENV, BEARER_LOGIN_PORT, SUCCESS_LOGIN_PAGE } from '../utils/constants'
import { askForString } from '../utils/prompts'

type Event = 'success' | 'error'

class LoginAction extends BaseAction {
  _server?: http.Server
  _verifier!: string
  _challenge!: string
  private _listerners: Record<Event, (() => void)[]> = { success: [], error: [] }

  async run() {
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
    this.logger.debug('authorize params %j', params)
    const url = `${this.logger.constants.LoginDomain}/authorize?${toParams(params)}`
    const spawned = await opn(url)
    await Promise.race([
      new Promise((resolve, reject) => {
        spawned.on('close', async (code: any, signal: any) => {
          if (code !== 0) {
            this.logger.warn(
              this.logger.colors.yellow(
                `Unable to open a browser. If you want to retrieve a token please follow these steps\n`
              )
            )
            this.logger.log(this.logger.colors.bold('1/ access the url below and follow the login process:\n\n'), url)
            this.logger.log()
            this.logger.log(
              this.logger.colors.bold(`2/ when you access the success page copy the token and paste it here`)
            )
            const token = await askForString('Token')
            await this.getToken(token)
            resolve()
          }
        })
      }),
      new Promise((resolve, reject) => {
        this.on('error', reject)
        this.on('success', resolve)
      })
    ])
    this.logger.debug('login done')
    cliUx.action.stop()
  }

  on = (event: Event, callback: () => void) => {
    this._listerners[event] = [...this._listerners[event], callback]
  }

  private stopServer = () => {
    return new Promise((resolve, _reject) => {
      this.logger.debug('stopping server')
      if (this._server) {
        this._server.close(() => {
          this.logger.debug('server stopped')
          resolve()
        })
      } else {
        resolve()
      }
    })
  }

  private startServer = async (): Promise<http.Server> => {
    // stop the server if successfully authenticated
    this.on('success', this.stopServer)
    this.on('error', this.stopServer)
    const port = await getPort({ port: BEARER_LOGIN_PORT })
    return new Promise((resolve, reject) => {
      if (port !== BEARER_LOGIN_PORT) {
        this.logger.error(`bearer login requires port ${BEARER_LOGIN_PORT} to be available`)
        reject()
      }
      this.logger.debug('starting server')
      const app = express()

      app.use((_req: express.Request, res: express.Response, next: express.NextFunction) => {
        res.setHeader('Connection', 'close')
        next()
      })

      app.get('/login/callback', (req: express.Request, res: express.Response) => {
        const code: string = req.query.code || ''
        try {
          res.send(SUCCESS_LOGIN_PAGE)
          setTimeout(() => {
            Promise.all([this.getToken(code), this.stopServer()]).then(() => {
              this.logger.debug('calling success listeners')
              this._listerners.success.map(cb => cb())
            })
          }, 0)
          res.end()
        } catch (e) {
          this.logger.debug(e)
          this.logger.error('Error while fetching token')
          this._listerners.error.map(cb => cb())
        }
      })

      const server = app.listen(BEARER_LOGIN_PORT, () => {
        this.logger.debug('server started')
        resolve(server)
      })
      server.addListener('connection', socket => {
        socket.setTimeout(0)
      })
    })
  }

  getToken = async (code: string) => {
    try {
      const { data: token } = await axios.post<TAccessToken>(`${this.logger.constants.LoginDomain}/oauth/token`, {
        code,
        grant_type: 'authorization_code',
        client_id: `${LOGIN_CLIENT_ID}`,
        code_verifier: `${this._verifier}`,
        redirect_uri: this.callbackUrl
      })

      this.logger.debug('saving token: %j', token)
      await this.logger.bearerConfig.storeToken(token)
      this.logger.success('Successfully logged in!! 🐻')
    } catch (e) {
      this.logger.error(e)
    }
  }

  get callbackUrl(): string {
    return process.env.BEARER_LOGIN_CALLBACK_URL || `${this.logger.constants.DeveloperPortalUrl}cli/login-callback`
  }
}

const loginFlow = createExport(LoginAction)

export default loginFlow

export async function promptToLogin(command: baseCommand) {
  command.log(command.colors.bold('⚠️ It looks like you are not logged in'))
  const { shoudlLogin } = await inquirer.prompt<{ shoudlLogin: boolean }>([
    {
      message: 'Would you like to login?',
      name: 'shoudlLogin',
      type: 'list',
      choices: [{ name: 'Yes', value: true }, { name: 'No', value: false }]
    }
  ])
  if (shoudlLogin) {
    await loginFlow(command)
  } else {
    command.exit(0)
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
