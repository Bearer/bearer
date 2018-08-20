import { flags } from '@oclif/command'

import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'

export default class GenerateIndex extends BaseCommand {
  static description = 'Generate Intent or Component'

  static flags = {
    help: flags.help({ char: 'h' })
  }

  static args = [{ name: 'name' }]

  @RequireScenarioFolder()
  async run() {
    const { args } = this.parse(GenerateIndex)
    this.log(`Generate component ${args.name}`)
  }
}
