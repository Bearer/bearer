import { flags } from '@oclif/command'

import BaseCommand from '../BaseCommand'
import { AuthType } from '../types'

const authTypes = [
  { name: 'OAuth2', value: AuthType.OAuth2 },
  { name: 'Basic Auth', value: AuthType.Basic },
  { name: 'API Key', value: AuthType.ApiKey },
  { name: 'NoAuth', value: AuthType.NoAuth }
]

export default class New extends BaseCommand {
  static description = 'Generate a new scenario'

  static flags = {
    help: flags.help({ char: 'h' }),
    authType: flags.string({
      char: 'a',
      description: 'Authorization type', // help description for flag
      options: authTypes.map(auth => auth.value) // only allow the value to be from a discrete set
    })
  }

  static args = [{ name: 'ScenarioName' }]

  async run() {
    const { args, flags } = this.parse(New)
    try {
      const name: string = args.ScenarioName || (await this.askForName())
      const authType: AuthType = (flags.authType as AuthType) || (await this.askForAuthType())
      await this.copyFiles(name, authType)
    } catch (e) {
      this.error('Could not generate new scenario' + e.toString())
    }
  }

  async copyFiles(name: string, authType: AuthType) {
    console.log(name, authType)
  }

  async askForName(): Promise<string> {
    const { name } = await this.inquirer.prompt<{ name: string }>([
      { message: 'Scenario name:', type: 'input', name: 'name' }
    ])
    return name.trim()
  }

  async askForAuthType(): Promise<AuthType> {
    const { authenticationType } = await this.inquirer.prompt<{ authenticationType: AuthType }>([
      {
        message: 'What kind of authentication do you want to use?',
        type: 'list',
        name: 'authenticationType',
        choices: authTypes
      }
    ])
    return authenticationType
  }
}
