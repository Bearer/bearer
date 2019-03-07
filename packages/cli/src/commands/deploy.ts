import BaseCommand from '../base-command'
import { RequireIntegrationFolder } from '../utils/decorators'

import PushCommand from './push'

export default class Deploy extends BaseCommand {
  static description = '[DEPRECATED] Deploys integration'
  static hidden = true

  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  @RequireIntegrationFolder()
  async run() {
    this.warn('This command is deprecated and will be removed soon.')
    this.warn(
      this.colors.bold('Please use this command instead: ') + this.colors.bold(this.colors.yellow('bearer push'))
    )
    try {
      await PushCommand.run(['--path', this.locator.integrationRoot])
      this.ux.action.stop()
    } catch (e) {
      this.error(e)
    }
  }
}
