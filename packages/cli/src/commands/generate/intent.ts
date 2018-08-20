import IntentType from '@bearer/types/lib/IntentTypes'
import { flags } from '@oclif/command'

import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'

const types = [
  { name: 'fetch', value: IntentType.FetchData },
  { name: 'save', value: IntentType.SaveState },
  { name: 'retrieve', value: IntentType.RetrieveState }
]

export default class GenerateIntent extends BaseCommand {
  static description = 'Generate a Bearer Intent'

  static flags = {
    help: flags.help({ char: 'h' }),
    type: flags.string({ char: 't', options: types.map(t => t.name) })
  }

  static args = [{ name: 'name' }]

  @RequireScenarioFolder()
  async run() {
    const { args, flags } = this.parse(GenerateIntent)
    this.log('Generating intent:')
    const type: IntentType = !flags.type ? await this.askForType() : types.find(t => t.name === flags.type)!.value
    const name = args.name || (await this.askForName())

    this.success(`Generated intent: name: ${name} type: ${type} `)
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
        choices: types.map(type => ({ ...type, name: this.case.pascal(type.name) }))
      }
    ])
    return type
  }
}
