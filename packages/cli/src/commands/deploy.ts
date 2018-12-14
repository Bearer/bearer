import BaseCommand from '../base-command'
import { RequireScenarioFolder } from '../utils/decorators'

import PushCommand from './push'

export default class Deploy extends BaseCommand {
  static description = '[DEPRECATED] Deploys scenario'
  static hidden = true

  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  @RequireScenarioFolder()
  async run() {
    this.warn('This command is deprecated and will be removed soon.')
    this.warn(
      this.colors.bold('Please use this command instead: ') + this.colors.bold(this.colors.yellow('bearer push'))
    )
    try {
      await PushCommand.run(['--path', this.locator.scenarioRoot])
      this.ux.action.stop()
    } catch (e) {
      this.error(e)
    }
  }
}
