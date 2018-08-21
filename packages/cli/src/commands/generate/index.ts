import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'

export default class GenerateIndex extends BaseCommand {
  static description = 'Generate Intent or Component'

  static flags = {
    ...BaseCommand.flags
  }

  static args = [{ name: 'name' }]

  @RequireScenarioFolder()
  async run() {
    const { args } = this.parse(GenerateIndex)
    // askForTYpe
    // AskForName
    this.log(`Generate component ${args.name}`)
  }
}
