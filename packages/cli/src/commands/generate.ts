import { flags } from '@oclif/command'

import BaseLegacyCommand from '../BaseLegacyCommand'

const blankComponent = 'blank-component'
const collectionComponent = 'collection-component'
const rootGroup = 'root-group'
const setup = 'setup'

export default class Generate extends BaseLegacyCommand {
  static description = 'Generate scenario intents or components'
  static flags = {
    help: flags.help({ char: 'h' }),
    [setup]: flags.boolean({ exclusive: [collectionComponent, blankComponent, rootGroup] }),
    [blankComponent]: flags.boolean({ exclusive: [setup, collectionComponent, rootGroup] }),
    [collectionComponent]: flags.boolean({ exclusive: [setup, blankComponent, rootGroup] }),
    [rootGroup]: flags.boolean({ exclusive: [setup, collectionComponent, blankComponent] })
  }

  static args = [{ name: 'name', required: false }]

  async run() {
    const { args, flags } = this.parse(Generate)
    const cmdArgs = []
    if (flags[blankComponent]) {
      cmdArgs.push(`--${blankComponent}`)
    } else if (flags[collectionComponent]) {
      cmdArgs.push(`--${collectionComponent}`)
    } else if (flags[rootGroup]) {
      cmdArgs.push(`--${rootGroup}`)
    } else if (flags[setup]) {
      cmdArgs.push(`--${setup}`)
    }
    this.debug('generate args', args)
    this.debug('generate flags', flags)
    if (args.name && !flags[setup]) {
      cmdArgs.push(args.name)
    }
    this.runLegacy(['generate', ...cmdArgs])
  }
}
