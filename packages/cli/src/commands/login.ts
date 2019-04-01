import { flags } from '@oclif/command'
import getPort from 'get-port'
import BaseCommand from '../base-command'
import * as http from 'http'

// const tokenUrl = 'https://app.bearer.sh/settings'

const BEARER_LOGIN_PORT = 56789

export default class Login extends BaseCommand {
  _server?: http.Server

  static description = 'Login to Bearer platform'

  static flags = {
    ...BaseCommand.flags,
    email: flags.string({ char: 'e' })
  }

  static args = []

  async run() {
    const port = await getPort({ port: BEARER_LOGIN_PORT })
    if (port !== BEARER_LOGIN_PORT) {
      this.error(`bearer login requires port ${BEARER_LOGIN_PORT} to be available`)
    }
    this.debug('starting server')
    this._server = http
      .createServer((request, response) => {
        let body = ''
        request.on('data', chunk => {
          body += chunk
        })
        request.on('end', () => {
          try {
            const data = JSON.parse(body)
            this.debug(data)
            response.setHeader('Access-Control-Allow-Origin', this.bearerConfig.DeveloperPortalUrl)
            response.write('OK')
            response.end()
            this.stopServer()
          } catch (e) {
            this.debug(e)
          }
        })
      })
      .listen(BEARER_LOGIN_PORT)

    // start server on specific port
    //   if not started raise
    // open browser
    // on response save code
    // ensure server is shutdown

    // const { flags } = this.parse(Login)
    // const email = process.env.BEARER_EMAIL || flags.email || (await this.askForString('Enter your email'))
    // const token = process.env.BEARER_TOKEN || (await this.askToken())
    // this.ux.action.start('Logging you in')
    // const status = await this.logUser(email, token)
    // this.ux.action.stop()
    // this.log(status)
  }

  private stopServer = () => {
    this.debug('stopping server')
    if (this._server) {
      this._server.close(() => {
        this.debug('server stopped')
      })
    }
  }

  // async logUser(username: string, accessToken: string): Promise<string> {
  //   try {
  //     const { statusCode, body } = await this.serviceClient.login({ Username: username, Password: accessToken })
  //     switch (statusCode as number) {
  //       case 200: {
  //         this.bearerConfig.storeBearerConfig({
  //           ...body.user,
  //           ExpiresAt: Date.now() + body.authorization.AuthenticationResult.ExpiresIn * 1000,
  //           authorization: body.authorization
  //         })
  //         return `Successfully logged in as ${username}! 🤘`
  //       }

  //       case 401: {
  //         this.error('Unauthorized: Invalid credentials')
  //         return 'There was an error while trying to login to bearer'
  //       }

  //       default: {
  //         this.error('Unhandled status')
  //         this.debug('status:', statusCode, 'body:', JSON.stringify(body))
  //         return 'An error occured. Please retry'
  //       }
  //     }
  //   } catch (e) {
  //     this.error(e)
  //     return 'An error occured'
  //   }
  // }

  // async askToken(): Promise<string> {
  //   this.log('')
  //   this.log(this.colors.italic(`Find your token at this location: ${this.colors.bold(tokenUrl)}`))
  //   const { token } = await this.inquirer.prompt<{ token: string }>([
  //     {
  //       message: 'Enter your token:',
  //       type: 'password',
  //       name: 'token'
  //     }
  //   ])
  //   return token
  // }
}
