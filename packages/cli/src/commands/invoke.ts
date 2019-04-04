import { flags } from '@oclif/command'

import BaseLegacyCommand from '../base-legacy-command'

export default class Invoke extends BaseLegacyCommand {
  static description = 'invoke Function locally'

  static flags = {
    help: flags.help({ char: 'h' }),
    path: flags.string({ char: 'p' })
  }

  static args = [{ name: 'Function_Name', required: true }]

  async run() {
    const { args, flags } = this.parse(Invoke)
    const cmdArgs = [args.Function_Name]
    if (flags.path) {
      cmdArgs.push(`--path=${flags.path}`)
    }
    this.runLegacy(['invoke', ...cmdArgs])
  }
}
