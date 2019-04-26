import Command, { flags } from '@oclif/command'
import * as colors from 'colors/safe'
import * as fs from 'fs'

import { AuthConfig, BaseConfig } from './types'
import Locator from './utils/locator'
import setupConfig, { Config } from './utils/setup-config'
import { devPortalClient } from './utils/devPortal'

export default abstract class extends Command {
  constants!: BaseConfig
  bearerConfig!: Config
  silent = false

  get locator() {
    return new Locator(this.bearerConfig)
  }

  get colors() {
    return colors
  }

  get integrationAuthConfig(): AuthConfig {
    return require(this.locator.authConfigPath)
  }

  get hasViews(): boolean {
    return fs.existsSync(this.locator.srcViewsDir)
  }

  static flags = {
    help: flags.help({ char: 'h' }),
    path: flags.string({}),
    silent: flags.boolean({})
    // logLevel: flags.string({ options: ['error', 'warn', 'info', 'debug'], default: 'info' })
  }

  success(message: string) {
    this.log(this.colors.green(message))
  }
  // @ts-ignore
  debug(...args: any[]) {
    super.debug.apply(this, args)
  }

  log(_message?: string, ..._args: any[]) {
    if (!this.silent) {
      // @ts-ignore
      super.log.apply(this, arguments)
    }
  }

  warn(_input: string | Error) {
    if (!this.silent) {
      // @ts-ignore
      super.warn.apply(this, arguments)
    }
  }

  get devPortalClient() {
    return devPortalClient(this)
  }

  get isIntegrationLocation() {
    return fs.existsSync(this.locator.authConfigPath)
  }

  // protected logLevel: any

  async init() {
    const { flags } = this.parse(this.constructor as any)
    const path = flags.path || undefined
    const { constants, config } = setupConfig(path)
    this.bearerConfig = config
    this.constants = constants
    this.silent = flags.silent
  }
}
