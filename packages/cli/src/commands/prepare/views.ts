import { flags } from '@oclif/command'
import * as path from 'path'

import BaseCommand from '../../base-command'
import { RequireScenarioFolder } from '../../utils/decorators'
import { copyFiles, ensureFolderExists, ensureSymlinked } from '../../utils/helpers'

export default class PrepareViews extends BaseCommand {
  static description = 'Prepare scenario views'
  static hidden = true
  static flags = {
    ...BaseCommand.flags,
    empty: flags.boolean()
  }

  static args = []

  @RequireScenarioFolder()
  async run() {
    const { flags } = this.parse(PrepareViews)
    // Prepare folder structure
    this.debug('Preparing views structure')
    ensureFolderExists(this.locator.buildViewsDir, flags.empty)
    ensureFolderExists(this.locator.buildViewsComponentsDir)

    this.debug('Symlinking node_modules')
    ensureSymlinked(
      this.locator.scenarioRootResourcePath('node_modules'),
      this.locator.buildViewsResourcePath('node_modules')
    )

    this.debug('Symlinking package.json')
    ensureSymlinked(
      this.locator.scenarioRootResourcePath('package.json'),
      this.locator.buildViewsResourcePath('package.json')
    )

    this.debug('Copying stencil config')
    const vars = { componentTagName: this.case.kebab(this.bearerConfig.scenarioConfig.scenarioTitle) }
    await copyFiles(this, 'start', this.locator.buildViewsDir, vars)

    ensureSymlinked(
      this.locator.buildViewsResourcePath('global'),
      path.join(this.locator.buildViewsComponentsDir, 'global')
    )
  }
}
