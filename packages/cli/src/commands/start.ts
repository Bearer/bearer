import { flags } from '@oclif/command'

import BaseLegacyCommand from '../BaseLegacyCommand'

import GenerateSetup from './generate/setup'
const noOpen = 'no-open'
const noInstall = 'no-install'
const noWatcher = 'no-watcher'

export default class Start extends BaseLegacyCommand {
  static description = 'Start local development environment'

  static flags = {
    help: flags.help({ char: 'h' }),
    [noOpen]: flags.boolean({}),
    [noInstall]: flags.boolean({}),
    [noWatcher]: flags.boolean({ hidden: true })
  }

  static args = []

  async run() {
    const { flags } = this.parse(Start)
    const cmdArgs = []
    if (flags[noOpen]) {
      cmdArgs.push(`--${noOpen}`)
    }
    if (flags[noInstall]) {
      cmdArgs.push(`--${noInstall}`)
    }
    if (flags[noWatcher]) {
      cmdArgs.push(`--${noWatcher}`)
    }

    await GenerateSetup.run(['--silent'])

    this.runLegacy(['start', ...cmdArgs])
  }
}
