import BaseCommand from '../../base-command'
import { RequireIntegrationFolder } from '../../utils/decorators'
import { ensureFolderExists } from '../../utils/helpers'

export default class PrepareIntents extends BaseCommand {
  static description = 'describe the command here'
  static hidden = true
  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  @RequireIntegrationFolder()
  async run() {
    ensureFolderExists(this.locator.buildIntentsDir, true)
    ensureFolderExists(this.locator.buildArtifactDir, true)
  }
}
