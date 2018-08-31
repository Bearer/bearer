import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'

export default class PackViews extends BaseCommand {
  static description = 'Pack scenario views'
  static hidden = true
  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  @RequireScenarioFolder()
  async run() {
    this.success('Packed views')
  }
}
