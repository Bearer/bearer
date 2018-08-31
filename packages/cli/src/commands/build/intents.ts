import * as Listr from 'listr'

import BaseCommand from '../../BaseCommand'
import installDependencies from '../../tasks/installDependencies'
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
    const tasks = new Listr([installDependencies({ cwd: this.locator.scenarioRoot })])
    try {
      await tasks.run()
      this.success('Built intents')
    } catch (e) {
      this.error(e)
    }
  }
}
