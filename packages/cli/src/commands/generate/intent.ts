import { flags } from '@oclif/command'

import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'

export default class GenerateIntent extends BaseCommand {
  static description = 'Generate a Bearer Intent'

  static flags = {
    help: flags.help({ char: 'h' })
  }

  static args = [{ name: 'name', required: true }]

  @RequireScenarioFolder()
  async run() {
    const { args } = this.parse(GenerateIntent)
    this.log(`Generate Intent ${args.name}`)
  }
}
