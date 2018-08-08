import { flags } from '@oclif/command'

import BaseLegacyCommand from '../BaseLegacyCommand'

export default class Link extends BaseLegacyCommand {
  static description = 'Link your local scenario to a remote one'

  static flags = {
    help: flags.help({ char: 'h' })
  }

  static args = [{ name: 'Scenario_Identifier', required: true }]

  async run() {
    const { args } = this.parse(Link)
    this.runLegacy(['link', args.Scenario_Identifier])
  }
}
