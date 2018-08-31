import { flags } from '@oclif/command'
import * as Listr from 'listr'

import BaseCommand from '../../BaseCommand'
import installDependencies from '../../tasks/installDependencies'
import { RequireScenarioFolder } from '../../utils/decorators'

const skipInstall = 'skip-install'

export default class BuildViews extends BaseCommand {
  static description = 'Build scenario views'
  static hidden = true
  static flags = {
    ...BaseCommand.flags,
    [skipInstall]: flags.boolean({})
  }

  static args = []

  @RequireScenarioFolder()
  async run() {
    const { flags } = this.parse(BuildViews)

    const tasks: Array<Listr.ListrTask> = [
      {
        title: 'Transpile views',
        task: async (_ctx: any, _task: any) => {
          this.transpile()
        }
      },
      {
        title: 'Build views',
        task: async (_ctx: any, _task: any) => {
          this.build()
        }
      }
    ]
    if (!flags[skipInstall]) {
      tasks.unshift(installDependencies({ cwd: this.locator.scenarioRoot }))
    }

    try {
      await new Listr(tasks).run()
      this.success('Built views')
    } catch (e) {
      this.error(e)
    }
  }

  transpile = async () => {}
  build = async () => {}
}
