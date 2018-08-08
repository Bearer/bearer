import { flags } from '@oclif/command'

import BaseLegacyCommand from '../BaseLegacyCommand'

const viewsOnly = 'views-only'
const intentsOnly = 'intents-only'

export default class Deploy extends BaseLegacyCommand {
  static description = 'Deploy a scenario'

  static flags = {
    help: flags.help({ char: 'h' }),
    [viewsOnly]: flags.boolean({ char: 's', exclusive: [intentsOnly], description: 'Deploy views only' }),
    [intentsOnly]: flags.boolean({ char: 'i', exclusive: [viewsOnly], description: 'Deploy intents only' })
  }

  static args = []

  async run() {
    const { flags } = this.parse(Deploy)
    const cmdArgs = []
    if (flags[viewsOnly]) {
      cmdArgs.push(`--${viewsOnly}`)
    } else if (flags[intentsOnly]) {
      cmdArgs.push(`--${intentsOnly}`)
    }
    this.runLegacy(['deploy', ...cmdArgs])
  }
}
