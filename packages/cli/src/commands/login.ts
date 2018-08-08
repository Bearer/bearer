import { flags } from '@oclif/command'

import BaseLegacyCommand from '../BaseLegacyCommand'

export default class Login extends BaseLegacyCommand {
  static description = 'Login to Bearer platform'

  static flags = {
    help: flags.help({ char: 'h' }),
    email: flags.string({ char: 'e', required: true })
  }

  static args = []

  async run() {
    const { flags } = this.parse(Login)
    const cmdArgs = []
    if (flags.email) {
      cmdArgs.push(`--email=${flags.email}`)
    }
    this.runLegacy(['login', ...cmdArgs])
  }
}
