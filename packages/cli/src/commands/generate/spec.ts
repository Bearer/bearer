import { flags } from '@oclif/command'
import * as fs from 'fs'
import * as path from 'path'

import BaseCommand from '../../base-command'
import { RequireIntegrationFolder, skipIfNoViews } from '../../utils/decorators'
import { copyFiles } from '../../utils/helpers'

export default class GenerateSpec extends BaseCommand {
  static description = 'Generate spec file for bearer integration'
  static hidden = true
  static flags = {
    ...BaseCommand.flags,
    force: flags.boolean({})
  }

  static args = []

  @skipIfNoViews()
  @RequireIntegrationFolder()
  async run() {
    const { flags } = this.parse(GenerateSpec)
    const targetFolder = this.locator.integrationRoot
    if (flags.force || !specExists(targetFolder)) {
      try {
        await copyFiles(this, 'generate/integration_specs', targetFolder, {})
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
