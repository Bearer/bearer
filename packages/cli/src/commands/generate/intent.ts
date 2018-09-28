import { templates } from '@bearer/templates'
import { Authentications } from '@bearer/types/lib/Authentications'
import IntentType from '@bearer/types/lib/IntentTypes'
import { flags } from '@oclif/command'
import * as inquirer from 'inquirer'

import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'
import { copyFiles } from '../../utils/helpers'

const types = [
  { name: 'Fetch', value: IntentType.FetchData, cli: 'fetch' },
  { name: 'Save State', value: IntentType.SaveState, cli: 'save' },
  { name: 'Retrieve Sate', value: IntentType.RetrieveState, cli: 'retrieve' }
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
    if (!templates[authType]) {
      // TODO: better error output
      this.error(
        `Incorrect AuthType please update "authType" field of auth.config.json within your scenario, with one of these values : ${Object.values(
          Authentications
        ).join('  |  ')}`
      )
    }
    try {
      const vars = this.getVars(name, type, authType)
      await copyFiles(this, `generate/intent`, this.locator.srcIntentsDir, vars)
      this.success(`\nIntent generated`)
      // TODO: add a nicer display
    } catch (e) {
      this.error(e)
    }
  }

  getVars(name: string, intentType: IntentType, authType: Authentications) {
    const actionExample = this.getActionExample(intentType, authType)
    return {
      fileName: name,
      intentName: name,
      intentClassName: this.case.pascal(name),
      authType,
      intentType,
      actionExample,
      callbackType: `T${intentType}Callback`
    }
  }

  getActionExample(intentType: IntentType, authType: Authentications): string {
    return templates[authType][intentType] as string
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
