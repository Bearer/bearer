import * as express from 'express'
import * as fs from 'fs'
import axios from 'axios'
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
import { startServer, TDestroyableServer, UnavailablePort } from '../../actions/startLocalServer'

type Event = 'success' | 'error'
const SEPARATOR = ':'
const COMMAND = `bearer setup:auth`
const enum keys {
  BEARER_AUTH_CLIENT_ID = 'BEARER_AUTH_CLIENT_ID',
  BEARER_AUTH_CLIENT_SECRET = 'BEARER_AUTH_CLIENT_SECRET',
  BEARER_AUTH_CONSUMER_KEY = 'BEARER_AUTH_CONSUMER_KEY',
  BEARER_AUTH_CONSUMER_SECRET = 'BEARER_AUTH_CONSUMER_SECRET',
  BEARER_AUTH_USERNAME = 'BEARER_AUTH_USERNAME',
  BEARER_AUTH_PASSWORD = 'BEARER_AUTH_PASSWORD',
  BEARER_AUTH_APIKEY = 'BEARER_AUTH_APIKEY'
}
export default class SetupAuth extends BaseCommand {
  static description = `setup API credentials for local development.
If you would like to bypass the prompt, you can either:
\t* pass credentials as argument (see description later)
\t* use environment variables
see examples`

  static flags = {
    ...BaseCommand.flags
  }

  static examples = [
    `With argument`,
    `\t${COMMAND} CLIENT_ID${SEPARATOR}CLIENT_SECRET`,
    `\t${COMMAND} CONSUMER_KEY${SEPARATOR}CONSUMER_SECRET`,
    `\t${COMMAND} USERNAME${SEPARATOR}PASSWORD`,
    `\t${COMMAND} APIKEY`,
    `With environment variables`,
    `\t${keys.BEARER_AUTH_CLIENT_ID}=CLIENT_ID ${keys.BEARER_AUTH_CLIENT_SECRET}=CLIENT_SECRET ${COMMAND}`,
    `\t${keys.BEARER_AUTH_CONSUMER_KEY}=CONSUMER_KEY ${keys.BEARER_AUTH_CONSUMER_SECRET}=CONSUMER_SECRET ${COMMAND}`,
    `\t${keys.BEARER_AUTH_USERNAME}=USERNAME ${keys.BEARER_AUTH_CONSUMER_SECRET}=PASSWORD ${COMMAND}`,
    `\t${keys.BEARER_AUTH_APIKEY}=APIKEY ${COMMAND}`
  ]

  static args = [
    {
      name: 'credentials',
      description: `Provide inline credentials`,
      required: false,
      default: ''
    }
  ]

  _server?: TDestroyableServer
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
    const { args } = this.parse(SetupAuth)
    const config: TConfig = JSON.parse(
      fs.readFileSync(this.locator.authConfigPath, {
        encoding: 'utf8'
      })
    )
    const { authType, provider, hint } = config as { authType: Authentications; provider?: string; hint?: string }

    if (hint) {
      this.log(this.colors.italic(this.colors.yellow(`\n${hint}\n`)))
    }

    function prefixedPrompt(message: string) {
      return prefix(message, provider)
    }

    switch (authType) {
      case Authentications.OAuth2: {
        const [idArg, secretArg] = args.credentials.split(SEPARATOR)
        const {
          [keys.BEARER_AUTH_CLIENT_ID]: id = idArg,
          [keys.BEARER_AUTH_CLIENT_SECRET]: secret = secretArg
        } = process.env
        const clientID = id || (await askForString(prefixedPrompt('Client ID'), { type: 'password' }))
        const clientSecret = secret || (await askForString(prefixedPrompt('Client secret'), { type: 'password' }))

        this.debug('Your credentials:\n%j', { ...config, clientID, clientSecret: clientSecret.replace(/./g, '*') })

        const newConfig = { ...config, clientID, clientSecret }
        const token = await this.fetchAuthToken(newConfig as configs.TOAuth2Config)
        const setup = JSON.parse(Buffer.from(token, 'base64').toString('ascii')) as contexts.OAuth2
        await this.persistSetup(setup)
        break
      }
      case Authentications.OAuth1: {
        const [keyArg, secretArg] = args.credentials.split(SEPARATOR)
        const {
          [keys.BEARER_AUTH_CONSUMER_KEY]: key = keyArg,
          [keys.BEARER_AUTH_CONSUMER_SECRET]: secret = secretArg
        } = process.env

        const consumerKey = key || (await askForString(prefixedPrompt('Consumer key')))
        const consumerSecret = secret || (await askForPassword(prefixedPrompt('Consumer secret')))

        this.debug('Your credentials:\n%j', { consumerKey, consumerSecret: consumerSecret.replace(/./g, '*') })

        const newConfig = { ...config, consumerKey, consumerSecret }
        const token = await this.fetchAuthToken(newConfig as configs.TOAuth1Config)
        const setup = JSON.parse(Buffer.from(token, 'base64').toString('ascii')) as contexts.OAuth1
        // TODO: fix this when context.OAuth1 is fixed  and well defined
        await this.persistSetup(setup)
        break
      }

      case Authentications.Basic: {
        const [usernameArg, passwordArg] = args.credentials.split(SEPARATOR)
        const {
          [keys.BEARER_AUTH_USERNAME]: basicUsername = usernameArg,
          [keys.BEARER_AUTH_PASSWORD]: basicPassword = passwordArg
        } = process.env
        const username = basicUsername || (await askForString(prefixedPrompt('Username')))
        const password = basicPassword || (await askForPassword(prefixedPrompt('Password')))

        await this.persistSetup({ username, password } as contexts.Basic)

        break
      }
      case Authentications.ApiKey: {
        const { [keys.BEARER_AUTH_APIKEY]: key } = process.env
        const apiKey = key || args.credentials || (await askForPassword(prefixedPrompt('API Key')))
        await this.persistSetup({ apiKey } as contexts.ApiKey)
        break
      }
      case Authentications.Custom:
      case Authentications.NoAuth: {
        await this.persistSetup({})
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

  async fetchAuthToken(config: configs.TOAuth2Config | configs.TOAuth1Config): Promise<TBase64EncodedString> {
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
          this.log(this.colors.bold('1/ access the url below and follow the login process:\n\n'), fallback)
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
    return new Promise((resolve, _reject) => {
      if (this._server) {
        this.debug('stopping server')
        this._server.destroy(() => {
          this.debug('server stopped')
          resolve()
        })
      } else {
        resolve()
      }
    })
  }

  private startServer = async (): Promise<TDestroyableServer> => {
    this.on('success', this.stopServer)
    this.on('error', this.stopServer)

    try {
      const app = express()

      app.use((_req: express.Request, res: express.Response, next: express.NextFunction) => {
        res.setHeader('Connection', 'close')
        next()
      })

      app.get('/setup/auth-callback', (req: express.Request, res: express.Response) => {
        const token: TBase64EncodedString = req.query.token || ''
        try {
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

      this.debug('starting server')
      return await startServer(BEARER_AUTH_PORT, app)
    } catch (e) {
      if (e instanceof UnavailablePort) {
        this.error('bearer setup:auth cannot start')
      }
      throw e
    }
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

function prefix(message: string, prefix?: string) {
  const prefixedMessage = [prefix, message].filter(Boolean).join(' ')

  return `Enter ${prefixedMessage}`
}

class UnreachableCaseError extends Error {
  constructor(val: never) {
    super(`Unreachable case: ${val}`)
  }
}

type TBase64EncodedString = string
