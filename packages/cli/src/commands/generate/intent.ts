import { Authentications } from '@bearer/types/lib/authentications'
import IntentType from '@bearer/types/lib/intent-types'
import { flags } from '@oclif/command'
import * as inquirer from 'inquirer'

import BaseCommand from '../../base-command'
import { RequireScenarioFolder } from '../../utils/decorators'
import generateIntent from '../../utils/templates/intents'

const types = [
  { name: 'Fetch', value: IntentType.FetchData, cli: 'fetch' },
  { name: 'Save State', value: IntentType.SaveState, cli: 'save' }
]

const typeChoices = [types.slice(0, 1)[0], new inquirer.Separator(), ...types.slice(1)]

export default class GenerateIntent extends BaseCommand {
  static description = 'Generate a Bearer Intent'
  static aliases = ['g:i']
  static flags = {
    ...BaseCommand.flags,
    type: flags.string({
      char: 't',
      options: types.map(t => t.cli)
    })
  }

  static args = [{ name: 'name' }]

  @RequireScenarioFolder()
  async run() {
    const { args, flags } = this.parse(GenerateIntent)
    const type: IntentType = !flags.type
      ? await this.askForType()
      : types.find(t => (t as { cli: string }).cli === flags.type)!.value
    const name = args.name || (await this.askForName())
    const authType = this.scenarioAuthConfig.authType

    if (!Object.values(Authentications).includes(authType)) {
      // TODO: better error output
      this.error(
        `Incorrect AuthType please update "authType" field of auth.config.json within your scenario, 
        with one of these values : ${Object.values(Authentications).join('  |  ')}`
      )
    }
    try {
      await generateIntent(this, authType, type, name)
      this.success(`\nIntent generated`)
    } catch (e) {
      this.error(e)
    }
  }

  async askForName(): Promise<string> {
    return this.askForString('Name')
  }

  async askForType(): Promise<IntentType> {
    const { type } = await this.inquirer.prompt<{ type: IntentType }>([
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
