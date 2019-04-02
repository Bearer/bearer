import Command from '@oclif/command'

import Locator from './utils/locator'
import setupConfig, { Config } from './utils/setup-config'
import * as fs from 'fs-extra'

export default abstract class extends Command {
  runLegacy(cmdArgs: any[]) {
    const cli = require('@bearer/bearer-cli/lib/bin/index').default
    this.debug('Legacy command arguments', JSON.stringify(cmdArgs))
    cli(['null', 'null'].concat(cmdArgs))
  }

  async init() {
    const path = process.cwd()
    this.bearerConfig = setupConfig(path).config
  }

  get locator() {
    return new Locator(this.bearerConfig)
  }

  get hasViews(): boolean {
    return fs.existsSync(this.locator.srcViewsDir)
  }

  bearerConfig!: Config
}
