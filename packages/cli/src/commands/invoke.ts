import { flags } from '@oclif/command'

import BaseLegacyCommand from '../base-legacy-command'

export default class Invoke extends BaseLegacyCommand {
  static description = 'Invoke Intent locally'

  static flags = {
    help: flags.help({ char: 'h' }),
    path: flags.string({ char: 'p' })
  }

  static args = [{ name: 'Intent_Name', required: true }]

  async run() {
    const { args, flags } = this.parse(Invoke)
    const cmdArgs = [args.Intent_Name]
    if (flags.path) {
      cmdArgs.push(`--path=${flags.path}`)
    }
    this.runLegacy(['invoke', ...cmdArgs])
  }
}
