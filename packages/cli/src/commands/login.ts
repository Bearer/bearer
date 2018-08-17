import { flags } from '@oclif/command'

import BaseCommand from '../BaseCommand'
const tokenUrl = 'https://app.bearer.sh/settings'

export default class Login extends BaseCommand {
  static description = 'Login to Bearer platform'

  static flags = {
    help: flags.help({ char: 'h' }),
    email: flags.string({ char: 'e' })
  }

  static args = []

  async run() {
    const { flags } = this.parse(Login)
    const email = flags.email || (await this.askEmail())
    const token = await this.askToken()
    this.ux.action.start('Logging you in')
    const status = await this.logUser(email, token)
    this.ux.action.stop(status)
  }

  async logUser(Username: string, AccessToken: string): Promise<string> {
    try {
      const { statusCode, body } = await this.serviceClient.login({ Username, Password: AccessToken })
      switch (statusCode as number) {
        case 200: {
          this.bearerConfig.storeBearerConfig({
            ...body.user,
            ExpiresAt: body.authorization.AuthenticationResult.ExpiresIn + Date.now(),
            authorization: body.authorization
          })
          return 'Successfully logged in! 🤘'
        }

        case 401: {
          this.error('Unauthorized')
          return 'There was an error while trying to login to bearer'
        }

        default: {
          this.error('Unhandled status')
          this.debug('status:', statusCode, 'body:', JSON.stringify(body))
          return 'An error occured. Please retry'
        }
      }
    } catch (e) {
      this.error(e)
      return 'An error occured'
    }
  }

  async askEmail(): Promise<string> {
    const { email } = await this.inquirer.prompt<{ email: string }>([
      {
        message: 'Enter your email:',
        name: 'email'
      }
    ])
    return email
  }

  async askToken(): Promise<string> {
    this.log(this.colors.italic(`Find your token at this location: ${this.colors.bold(tokenUrl)}`))
    this.log('\n')
    const { password } = await this.inquirer.prompt<{ password: string }>([
      {
        message: 'Enter your token:',
        type: 'password',
        name: 'password'
      }
    ])
    return password
  }
}
