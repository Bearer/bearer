import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'

export default class PackIntents extends BaseCommand {
  static description = 'Pack scenario intents'
  static hidden = true
  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  @RequireScenarioFolder()
  async run() {
    this.success('Packed intents')
  }
}
