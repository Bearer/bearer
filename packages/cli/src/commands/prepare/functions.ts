import BaseCommand from '../../base-command'
import { RequireIntegrationFolder } from '../../utils/decorators'
import { ensureFolderExists } from '../../utils/helpers'

export default class PrepareFunctions extends BaseCommand {
  static description = 'describe the command here'
  // TODO: remove prepare:intents when people have migrated to function
  static aliases = ['prepare:intents', 'p:f']
  static hidden = true
  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  @RequireIntegrationFolder()
  async run() {
    ensureFolderExists(this.locator.buildFunctionsDir, true)
    ensureFolderExists(this.locator.buildArtifactDir, true)
  }
}
