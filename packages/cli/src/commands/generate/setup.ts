import { flags } from '@oclif/command'
import * as fs from 'fs'
import * as path from 'path'

import BaseCommand from '../../base-command'
import { RequireIntegrationFolder, skipIfNoViews } from '../../utils/decorators'
import * as Listr from 'listr'
import buildSetup from '../../tasks/build-setup'

export default class GenerateSetup extends BaseCommand {
  static description = 'Generate a Bearer Setup'
  static hidden = true
  static flags = { ...BaseCommand.flags, force: flags.boolean({}) }

  static args = []

  @skipIfNoViews()
  @RequireIntegrationFolder()
  async run() {
    const { flags } = this.parse(GenerateSetup)
    if (flags.force || !setupExists(this.locator.srcViewsDir)) {
      const fields = this.integrationAuthConfig.setupViews
      if (fields && fields.length) {
        try {
          const vars = this.getVars(this.bearerConfig.integrationConfig.integrationTitle, fields)
          const tasks: Listr.ListrTask[] = buildSetup({
            vars,
            cmd: this
          })
          await new Listr(tasks).run()
        } catch (e) {
          this.error(e)
        }
      } else {
        this.warn(`No setupViews key found within auth.config.json file. skipping`)
      }
    } else {
      this.warn('Setup files already exist, use --force flag to overwrite existing ones.')
    }
  }

  getVars(integrationTitle: string, fields: any) {
    return {
      componentName: this.case.pascal(integrationTitle),
      componentTagName: this.case.kebab(integrationTitle),
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
