import { flags } from '@oclif/command'

import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'

export default class GenerateSetup extends BaseCommand {
  static description = 'Generate a Bearer Setup'
  static hidden = true
  static flags = {
    help: flags.help({ char: 'h' })
  }

  static args = [{ name: 'name', required: true }]

  @RequireScenarioFolder()
  async run() {
    const { args } = this.parse(GenerateSetup)
    this.log(`Generate Setup ${args.name}`)
  }
}
