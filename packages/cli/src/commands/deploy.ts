import { flags } from '@oclif/command'

import BaseLegacyCommand from '../BaseLegacyCommand'
import { ensureFreshToken } from '../utils/decorators'

import GenerateSpec from './generate/spec'
import PrepareViews from './prepare/views'

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

  @ensureFreshToken()
  async run() {
    const { flags } = this.parse(Deploy)
    const cmdArgs = []
    if (flags[viewsOnly]) {
      cmdArgs.push(`--${viewsOnly}`)
    } else if (flags[intentsOnly]) {
      cmdArgs.push(`--${intentsOnly}`)
    }
    await GenerateSpec.run(['--silent'])

    await PrepareViews.run(['--silent', '--empty'])

    this.runLegacy(['deploy', ...cmdArgs])
  }
}
