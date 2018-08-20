import { flags } from '@oclif/command'

import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'

export default class GenerateComponent extends BaseCommand {
  static description = 'Generate a Bearer component'

  static flags = {
    help: flags.help({ char: 'h' })
  }

  static args = [{ name: 'name', required: true }]

  @RequireScenarioFolder()
  async run() {
    const { args } = this.parse(GenerateComponent)
    this.log(`Generate component ${args.name}`)
  }
}
