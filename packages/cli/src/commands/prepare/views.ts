import { flags } from '@oclif/command'
import * as path from 'path'
import * as Case from 'case'

import BaseCommand from '../../base-command'
import { RequireIntegrationFolder, skipIfNoViews } from '../../utils/decorators'
import { copyFiles, ensureFolderExists, ensureSymlinked } from '../../utils/helpers'

export default class PrepareViews extends BaseCommand {
  static description = 'Prepare integration views'
  static hidden = true
  static flags = {
    ...BaseCommand.flags,
    empty: flags.boolean()
  }

  static args = []

  @skipIfNoViews()
  @RequireIntegrationFolder()
  async run() {
    const { flags } = this.parse(PrepareViews)
    // Prepare folder structure
    this.debug('Preparing views structure')
    ensureFolderExists(this.locator.buildViewsDir, flags.empty)
    ensureFolderExists(this.locator.buildViewsComponentsDir)

    this.debug('Symlinking node_modules')
    ensureSymlinked(
      this.locator.integrationRootResourcePath('node_modules'),
      this.locator.buildViewsResourcePath('node_modules')
    )

    this.debug('Symlinking package.json')
    ensureSymlinked(
      this.locator.integrationRootResourcePath('package.json'),
      this.locator.buildViewsResourcePath('package.json')
    )

    this.debug('Copying stencil config')
    const vars = { componentTagName: Case.kebab(this.bearerConfig.integrationConfig.integrationTitle) }
    await copyFiles(this, 'start', this.locator.buildViewsDir, vars)

    ensureSymlinked(
      this.locator.buildViewsResourcePath('global'),
      path.join(this.locator.buildViewsComponentsDir, 'global')
    )
  }
}
