import BaseCommand from '../BaseCommand'
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
    this.ux.action.start('Running bearer push')
    try {
      await PushCommand.run(['--path', this.locator.scenarioRoot])
      this.ux.action.stop()
      this.success(`🐻 Scenario successfully pushed.\n`)
      this.log(
        `Your scenario will be available soon at this location: ` +
          this.colors.bold(`${this.bearerConfig.DeveloperPortalUrl}scenarios/${this.bearerConfig.scenarioUuid}/preview`)
      )
    } catch (e) {
      this.error(e)
    }
  }
}
