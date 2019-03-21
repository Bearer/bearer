import { Authentications } from '@bearer/types/lib/authentications'
import FunctionType from '@bearer/types/lib/function-types'

import BaseCommand from '../../base-command'
import { RequireIntegrationFolder } from '../../utils/decorators'
import generateFunction from '../../utils/templates/functions'

export default class GenerateFunction extends BaseCommand {
  static description = 'Generate a Bearer Function'
  static aliases = ['g:f']
  static flags = { ...BaseCommand.flags }

  static args = [{ name: 'name' }]

  @RequireIntegrationFolder()
  async run() {
    const { args } = this.parse(GenerateFunction)
    const name = args.name || (await this.askForName())
    const authType = this.integrationAuthConfig.authType

    if (!Object.values(Authentications).includes(authType)) {
      // TODO: better error output
      this.error(
        `Incorrect AuthType please update "authType" field of auth.config.json within your integration, 
        with one of these values : ${Object.values(Authentications).join('  |  ')}`
      )
    }
    try {
      await generateFunction(this, authType, FunctionType.FetchData, name)
      this.success(`\nFunction generated`)
    } catch (e) {
      this.error(e)
    }
  }

  async askForName(): Promise<string> {
    return this.askForString('Name')
  }
}
