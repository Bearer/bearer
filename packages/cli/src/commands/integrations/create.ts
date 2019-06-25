import { flags } from '@oclif/command'
// @ts-ignore
import * as suggest from 'inquirer-prompt-suggest'

import BaseCommand from '../../base-command'
import CreateIntegration from '../../actions/createIntegration'

export default class IntegrationsCreate extends BaseCommand {
  static description = 'create a new Integration'

  static flags = {
    ...BaseCommand.flags,
    description: flags.string({ char: 'd' }),
    name: flags.string({ char: 'n' }),
    skipLink: flags.boolean({ char: 'l' })
  }

  async run() {
    const {
      flags: { name, description, skipLink }
    } = this.parse(IntegrationsCreate)
    await CreateIntegration(this, { name, description, link: Boolean(!skipLink) })
  }
}
