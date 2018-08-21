import * as serviceClient from '@bearer/bearer-cli/dist/src/lib/serviceClient'
import Command, { flags } from '@oclif/command'
import * as Case from 'case'
import cliUx from 'cli-ux'
import * as colors from 'colors/safe'
import * as copy from 'copy-template-dir'
import * as inquirer from 'inquirer'

import { AuthConfig, Config } from './types'
import Locator from './utils/locator'
import setupConfig from './utils/setupConfig'

export default abstract class extends Command {
  get locator() {
    return new Locator(this.bearerConfig)
  }

  get inquirer() {
    return inquirer
  }

  get copy() {
    return copy
  }

  get case() {
    return Case
  }

  get ux() {
    return cliUx
  }

  get colors() {
    return colors
  }

  get serviceClient() {
    return serviceClient(this.bearerConfig.IntegrationServiceUrl)
  }

  static flags = {
    help: flags.help({ char: 'h' }),
    path: flags.string({})
    // logLevel: flags.string({ options: ['error', 'warn', 'info', 'debug'], default: 'info' })
  }

  bearerConfig!: Config

  success(message: string) {
    this.log(this.colors.green(message))
  }

  // protected logLevel: any

  async init() {
    const { flags } = this.parse(this.constructor as any)
    const path = flags.path || undefined
    this.bearerConfig = setupConfig(path)
  }

  get scenarioAuthConfig(): AuthConfig {
    return require(this.locator.authConfigPath)
  }

  /**
   * Interactivity helpers
   */

  protected async askForString(message: string): Promise<string> {
    const { string } = await this.inquirer.prompt<{ string: string }>([
      {
        message: `${message}:`,
        name: 'string'
      }
    ])
    return string
  }
}
