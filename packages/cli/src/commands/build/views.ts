import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'

export default class BuildViews extends BaseCommand {
  static description = 'Build scenario views'
  static hidden = true
  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  @RequireScenarioFolder()
  async run() {
    this.success('Built views')
  }
}
