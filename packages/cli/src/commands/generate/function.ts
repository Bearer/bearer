import { Authentications } from '@bearer/types/lib/authentications'
import FunctionType from '@bearer/types/lib/function-types'
import { flags } from '@oclif/command'
import * as inquirer from 'inquirer'

import BaseCommand from '../../base-command'
import { RequireIntegrationFolder } from '../../utils/decorators'
import generateFunction from '../../utils/templates/functions'

const types = [
  { name: 'Fetch', value: FunctionType.FetchData, cli: 'fetch' },
  { name: 'Save State', value: FunctionType.SaveState, cli: 'save' }
]

const typeChoices = [types.slice(0, 1)[0], new inquirer.Separator(), ...types.slice(1)]

export default class GenerateFunction extends BaseCommand {
  static description = 'Generate a Bearer Function'
  static aliases = ['g:f']
  static flags = {
    ...BaseCommand.flags,
    type: flags.string({
      char: 't',
      options: types.map(t => t.cli)
    })
  }

  static args = [{ name: 'name' }]

  @RequireIntegrationFolder()
  async run() {
    const { args, flags } = this.parse(GenerateFunction)
    const type: FunctionType = !flags.type
      ? await this.askForType()
      : types.find(t => (t as { cli: string }).cli === flags.type)!.value
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
      await generateFunction(this, authType, type, name)
      this.success(`\nFunction generated`)
    } catch (e) {
      this.error(e)
    }
  }

  async askForName(): Promise<string> {
    return this.askForString('Name')
  }

  async askForType(): Promise<FunctionType> {
    const { type } = await this.inquirer.prompt<{ type: FunctionType }>([
      {
        message: 'Type:',
        type: 'list',
        name: 'type',
        choices: typeChoices
      }
    ])
    return type
  }
}
