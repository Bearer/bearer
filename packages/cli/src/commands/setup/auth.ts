import * as express from 'express'
import * as fs from 'fs'
import * as http from 'http'
import axios from 'axios'
import getPort from 'get-port'
import { modify, applyEdits } from 'jsonc-parser'
import { TConfig, Authentications, configs } from '@bearer/types/lib/authentications'
import { contexts } from '@bearer/functions/lib/declaration'

// TODO: use esModuleInterop config
// @ts-ignore
import * as open from 'open'

import BaseCommand from '../../base-command'
import { RequireIntegrationFolder } from '../../utils/decorators'
import { BEARER_AUTH_PORT, SUCCESS_LOGIN_PAGE } from '../../utils/constants'
import { askForString, askForPassword } from '../../utils/prompts'

type Event = 'success' | 'error'

export default class SetupAuth extends BaseCommand {
  static description = 'setup API credentials for local development'

  static flags = {
    ...BaseCommand.flags
  }

  _server?: http.Server
  _verifier!: string
  _challenge!: string
  private _listerners: {
    success: ((data: TBase64EncodedString) => void)[]
    error: ((data: any) => void)[]
  } = {
    success: [],
    error: []
  }

  @RequireIntegrationFolder()
  async run() {
    const config: TConfig = JSON.parse(
      fs.readFileSync(this.locator.authConfigPath, {
        encoding: 'utf8'
      })
    )
    const { authType } = config
    switch (authType) {
      case Authentications.OAuth2: {
        const { BEARER_AUTH_CLIENT_ID, BEARER_AUTH_CLIENT_SECRET } = process.env
        const clientID = BEARER_AUTH_CLIENT_ID || (await askForString('Client ID', { type: 'password' }))
        const clientSecret = BEARER_AUTH_CLIENT_SECRET || (await askForString('Client secret', { type: 'password' }))
        const newConfig = { ...config, clientID, clientSecret }
        this.debug('Your credentials:\n%j', { ...config, clientID, clientSecret: clientSecret.replace(/./g, '*') })
        const token = await this.fetchAuthToken(newConfig as configs.TOAuth2Config)
        const setup = JSON.parse(Buffer.from(token, 'base64').toString('ascii')) as contexts.OAuth2
        await this.persistSetup(setup)
        break
      }
      case Authentications.OAuth1: {
        const consumerKey = process.env.BEARER_AUTH_CONSUMER_KEY || (await askForString('Consumer key'))
        const consumerSecret = process.env.BEARER_AUTH_CONSUMER_SECRET || (await askForPassword('Consumer secret'))
        this.debug('Your credentials:\n%j', { consumerKey, consumerSecret: consumerSecret.replace(/./g, '*') })
        const newConfig = { ...config, consumerKey, consumerSecret }
        const token = await this.fetchAuthToken(newConfig as configs.TOAuth1Config)
        const setup = JSON.parse(Buffer.from(token, 'base64').toString('ascii')) as contexts.OAuth1
        // TODO: fix this when context.OAuth1 is fixed  and well defined
        await this.persistSetup(setup)
        break
      }

      case Authentications.Basic: {
        const username = await askForString('Username')
        const password = await askForPassword('Password')
        await this.persistSetup({ username, password } as contexts.Basic)
        break
      }
      case Authentications.ApiKey: {
        const apiKey = await askForPassword('API Key')
        await this.persistSetup({ apiKey } as contexts.ApiKey)
        break
      }
      case Authentications.Custom:
      case Authentications.NoAuth: {
        return this.warn(
          `The current authentication type of this integration is not supported by this command: ${authType}`
        )
      }
      default:
        // unsure we handled all authentications
        // http://ideasintosoftware.com/exhaustive-switch-in-typescript/
        this.error(`The current authentication type of this integration is not supported by this command: ${authType}`)
        throw new UnreachableCaseError(authType)
    }
  }

  fetchAuthToken = async (config: configs.TOAuth2Config | configs.TOAuth1Config): Promise<TBase64EncodedString> => {
    return new Promise(async (resolve, reject) => {
      this._server = await this.startServer()
      const redirectLocation = await axios
        .post(`${this.constants.IntegrationServiceHost}v2/auth/local-auth`, { config }, { maxRedirects: 0 })
        .catch(e => e.response.headers.location)

      // add listeners so that we can handle redirect and fallback the same way
      this.on<TBase64EncodedString>('success', data => {
        resolve(data)
      })

      this.on('error', data => {
        this.debug(data)
        reject('Error while receiving token')
      })

      const { url, fallback } = getOpeningUrls(
        `${this.constants.IntegrationServiceHost}v2/auth/${redirectLocation.replace('./', '')}&clientId=NONE`
      )
      this.debug('config: %j, location: %s', config, url)
      const a = await open(url)
      a.on('close', (code: any, signal: any) => {
        if (code !== 0) {
          // we are no able to open a browser so we ask the user to fill with the token
          this.stopServer()
          this.warn(
            this.colors.yellow(`Unable to open a browser. If you want to retrieve a token please follow these steps\n`)
          )
          this.log(this.colors.bold('1/ access the url below  and follow the login process:\n\n'), fallback)
          this.log()
          this.log(this.colors.bold(`2/ when you access the success page copy the token and paste it here`))
          askForString('Token').then((token: TBase64EncodedString) => {
            resolve(token)
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
      })
    }
  }

  private startServer = async (): Promise<http.Server> => {
    this.on('success', this.stopServer)
    this.on('error', this.stopServer)

    const port = await getPort({ port: BEARER_AUTH_PORT })

    return new Promise((resolve, reject) => {
      if (port !== BEARER_AUTH_PORT) {
        this.error(`bearer setup requires port ${BEARER_AUTH_PORT} to be available`)
        reject()
      }
      this.debug('starting server')

      const app = express()
      app.use((_req: express.Request, res: express.Response, next: express.NextFunction) => {
        res.setHeader('Connection', 'close')
        next()
      })
      app.get('/setup/auth-callback', (req: express.Request, res: express.Response) => {
        const token: TBase64EncodedString = req.query.token || ''
        try {
          res.setHeader('Connection', 'close')
          res.send(SUCCESS_LOGIN_PAGE).end()
          setTimeout(() => {
            this._listerners.success.map(cb => cb(token))
          }, 0)
        } catch (e) {
          this.debug(e)
          this._listerners.error.map(cb => cb(e))
          this.error('Error while handling callback')
        }
      })

      const server = app.listen(BEARER_AUTH_PORT, () => {
        this.debug('server started')
        resolve(server)
      })
    })
  }

  private persistSetup(config: contexts.OAuth2 | contexts.OAuth1 | contexts.ApiKey | contexts.Basic) {
    const setupRc = this.locator.localConfigPath
    if (!fs.existsSync(setupRc)) {
      fs.writeFileSync(setupRc, '{}', { encoding: 'utf8' })
    }
    const rawSetup = fs.readFileSync(setupRc, { encoding: 'utf8' })
    const updates = modify(rawSetup, ['setup', 'auth'], config, {
      formattingOptions: {
        insertSpaces: true,
        tabSize: 2,
        eol: '\n'
      }
    })
    const newSetupRc = applyEdits(rawSetup, updates)
    this.debug('Writing setup config\n%j', { config, setupRc: newSetupRc })
    fs.writeFileSync(setupRc, newSetupRc, { encoding: 'utf8' })
    this.success(`Auth credentials have been saved to ${this.locator.toRelative(setupRc)}`)
  }

  private on = <T>(event: Event, callback: (data: T) => void) => {
    this._listerners[event] = [...this._listerners[event], callback as any]
  }
}

function getOpeningUrls(url: string) {
  return {
    url: `${url}&localHostRedirectSupported=true`,
    fallback: `${url}&localHostRedirectSupported=inline`
  }
}
class UnreachableCaseError extends Error {
  constructor(val: never) {
    super(`Unreachable case: ${val}`)
  }
}

type TBase64EncodedString = string
