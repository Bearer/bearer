import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'

export default class BuildIntents extends BaseCommand {
  static description = 'Build scenario intents'
  static hidden = true
  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  @RequireScenarioFolder()
  async run() {
    this.success('Built intents')
  }
}
