import * as express from 'express'
import * as crypto from 'crypto'
import axios from 'axios'
import * as inquirer from 'inquirer'

import baseCommand from '../base-command'
import { startServer, TDestroyableServer, UnavailablePort } from '../actions/startLocalServer'

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
  _server?: TDestroyableServer
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
    this.logger.success('Successfully logged in!! 🐻')
    cliUx.action.stop()
  }

  on = (event: Event, callback: () => void) => {
    this._listerners[event] = [...this._listerners[event], callback]
  }

  private stopServer = () =>
    new Promise((resolve, _reject) => {
      if (this._server) {
        this.logger.debug('stopping server')
        this._server.destroy(() => {
          this.logger.debug('server stopped')
          resolve()
        })
      } else {
        resolve()
      }
    })

  private startServer = async (): Promise<TDestroyableServer> => {
    this.on('error', this.stopServer)
    try {
      const app = express()

      app.use((_req: express.Request, res: express.Response, next: express.NextFunction) => {
        res.setHeader('Connection', 'close')
        next()
      })
      app.get('/login/callback', (req: express.Request, res: express.Response, next: express.NextFunction) => {
        const code: string = req.query.code || ''
        try {
          res.send(SUCCESS_LOGIN_PAGE).end()
          Promise.all([this.getToken(code), this.stopServer()]).then(() => {
            this.logger.debug('calling success listeners')
            this._listerners.success.map(cb => cb())
          })
        } catch (e) {
          this.logger.debug(e)
          this.logger.error('Error while fetching token')
          this._listerners.error.map(cb => cb())
        }
      })

      this.logger.debug('starting server')
      return await startServer(BEARER_LOGIN_PORT, app)
    } catch (e) {
      if (e instanceof UnavailablePort) {
        this.logger.error('bearer login cannot start')
      }
      throw e
    }
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
  const { shouldLogin } = await inquirer.prompt<{ shouldLogin: boolean }>([
    {
      message: 'Would you like to login?',
      name: 'shouldLogin',
      type: 'list',
      choices: [{ name: 'Yes', value: true }, { name: 'No', value: false }]
    }
  ])
  if (shouldLogin) {
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
