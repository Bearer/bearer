import { flags } from '@oclif/command'

import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'

enum TComponent {
  BLANK = 'blank',
  COLLECTION = 'collection',
  ROOT = 'root'
}

export default class GenerateComponent extends BaseCommand {
  static description = 'Generate a Bearer component'

  static flags = {
    help: flags.help({ char: 'h' }),
    type: flags.string({ char: 't', options: Object.values(TComponent) })
  }

  static args = [{ name: 'name', required: true }]

  @RequireScenarioFolder()
  async run() {
    const { args, flags } = this.parse(GenerateComponent)
    const type: TComponent = (flags.type as TComponent) || (await this.askForComponentType())
    const name: string = args.name || (await this.askForName())
    // copy files
    // display create (same as init)
    this.success(`Generated component: name: ${name} | type: ${type}`)
  }

  async askForComponentType(): Promise<TComponent> {
    const { type } = await this.inquirer.prompt<{ type: TComponent }>([
      {
        message: 'What kind of component would you like to generate:',
        type: 'list',
        name: 'type',
        choices
      }
    ])
    return type
  }

  async askForName(): Promise<string> {
    return this.askForString('Name')
  }
}

// TODO: better names
const choices = [
  {
    name: 'Blank',
    value: TComponent.BLANK
  },
  {
    name: 'Root component',
    value: TComponent.ROOT
  },
  {
    name: 'Collection',
    value: TComponent.COLLECTION
  }
]
