import serviceClient from '@bearer/bearer-cli/lib/lib/serviceClient'
import Command, { flags } from '@oclif/command'
import * as Case from 'case'
import cliUx from 'cli-ux'
import * as colors from 'colors/safe'
import * as copy from 'copy-template-dir'
import * as inquirer from 'inquirer'

import { AuthConfig, Config } from './types'
import Locator from './utils/locator'
import scenarioClientFactory, { ScenarioClient } from './utils/scenario-client'
import setupConfig from './utils/setup-config'

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

  // TODO: fix typing
  get ux(): any {
    return cliUx
  }

  get colors() {
    return colors
  }

  // TODO: fix typing
  get serviceClient(): any {
    return serviceClient(this.bearerConfig.IntegrationServiceUrl)
  }

  get scenarioClient(): ScenarioClient {
    return scenarioClientFactory(this)
  }

  get scenarioAuthConfig(): AuthConfig {
    return require(this.locator.authConfigPath)
  }

  static flags = {
    help: flags.help({ char: 'h' }),
    path: flags.string({}),
    silent: flags.boolean({})
    // logLevel: flags.string({ options: ['error', 'warn', 'info', 'debug'], default: 'info' })
  }

  bearerConfig!: Config
  silent = false

  success(message: string) {
    this.log(this.colors.green(message))
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

  // protected logLevel: any

  async init() {
    const { flags } = this.parse(this.constructor as any)
    const path = flags.path || undefined
    this.silent = flags.silent
    this.bearerConfig = setupConfig(path)
  }

  /**
   * Interactivity helpers
   */

  protected async askForString(phrase: string): Promise<string> {
    const { response } = await this.inquirer.prompt<{ response: string }>([
      {
        message: `${phrase}:`,
        name: 'response'
      }
    ])
    return response
  }
}
