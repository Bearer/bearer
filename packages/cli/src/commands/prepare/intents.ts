import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'
import { ensureFolderExists } from '../../utils/helpers'

export default class PrepareIntents extends BaseCommand {
  static description = 'describe the command here'
  static hidden = true
  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  @RequireScenarioFolder()
  async run() {
    ensureFolderExists(this.locator.buildIntentsDir, true)
    ensureFolderExists(this.locator.intentsArtifactDir, true)
  }
}
