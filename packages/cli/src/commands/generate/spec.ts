import { flags } from '@oclif/command'
import * as fs from 'fs'
import * as path from 'path'

import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'
import { copyFiles } from '../../utils/helpers'

export default class GenerateSpec extends BaseCommand {
  static description = 'Generate spec file for bearer scenario'
  static hidden = true
  static flags = {
    ...BaseCommand.flags,
    force: flags.boolean({})
  }

  static args = []

  @RequireScenarioFolder()
  async run() {
    const { flags } = this.parse(GenerateSpec)
    const targetFolder = this.locator.scenarioRoot
    if (flags.force || !specExists(targetFolder)) {
      try {
        await copyFiles(this, 'generate/scenario_specs', targetFolder, {})
        this.success('Spec file successfully generated! 🎉')
      } catch (e) {
        this.error(e)
      }
    } else {
      this.warn('Spec file already exists, use --force flag to overwrite existing one.')
    }
  }
}

// Note: using Or condition in case the developer delete one but customized the other component
function specExists(location: string): boolean {
  return fs.existsSync(path.join(location, 'spec.ts'))
}
