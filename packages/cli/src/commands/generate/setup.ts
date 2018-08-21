import { flags } from '@oclif/command'
import * as fs from 'fs'
import * as path from 'path'

import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'
import { copyFiles } from '../../utils/helpers'

export default class GenerateSetup extends BaseCommand {
  static description = 'Generate a Bearer Setup'
  static hidden = true
  static flags = {
    ...BaseCommand.flags,
    help: flags.help({ char: 'h' }),
    force: flags.boolean({})
  }

  static args = []

  @RequireScenarioFolder()
  async run() {
    const { flags } = this.parse(GenerateSetup)
    if (flags.force || !setupExists(this.locator.srcViewsDir)) {
      const fields = this.scenarioAuthConfig.setupViews
      if (fields && fields.length) {
        try {
          await copyFiles(
            this,
            'generate/setup',
            this.locator.srcViewsDir,
            this.getVars(this.bearerConfig.scenarioConfig.scenarioTitle, fields)
          )
          this.success('Setup components successfully generated! 🎉')
        } catch (e) {
          this.error(e)
        }
      } else {
        this.warn(`No setupViews key found within auth.config.json file. skipping`)
      }
    } else {
      this.warn('Setup files already exists, use --force flag to overwrite existing ones.')
    }
  }

  getVars(scenarioTitle: string, fields: any) {
    return {
      componentName: this.case.pascal(scenarioTitle),
      componentTagName: this.case.kebab(scenarioTitle),
      fields: JSON.stringify(fields)
    }
  }
}

// Note: using Or condition in case the developer delete one but customized the other component
function setupExists(location: string): boolean {
  return (
    fs.existsSync(path.join(location, 'setup-action.tsx')) || fs.existsSync(path.join(location, 'setup-display.tsx'))
  )
}
