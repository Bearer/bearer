import { flags } from '@oclif/command'

import BaseLegacyCommand from '../BaseLegacyCommand'

export default class New extends BaseLegacyCommand {
  static description = 'Generate a new scenario'

  static flags = {
    help: flags.help({ char: 'h' })
  }

  static args = [{ name: 'ScenarioName', required: true }]

  async run() {
    const { args } = this.parse(New)
    this.runLegacy(['new', args.ScenarioName])
  }
}
