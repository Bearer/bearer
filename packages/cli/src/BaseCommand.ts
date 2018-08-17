import * as serviceClient from '@bearer/bearer-cli/dist/src/lib/serviceClient'
import Command from '@oclif/command'
import * as Case from 'case'
import cliUx from 'cli-ux'
import * as colors from 'colors/safe'
import * as copy from 'copy-template-dir'
import * as inquirer from 'inquirer'

import { Config } from './types'
import Locator from './utlis/locator'
import setupConfig from './utlis/setupConfig'

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
    // logLevel: flags.string({ options: ['error', 'warn', 'info', 'debug'], default: 'info' })
  }

  protected bearerConfig!: Config

  success(message: string) {
    this.log(this.colors.green(message))
  }

  // protected logLevel: any

  async init() {
    this.bearerConfig = setupConfig()
    // const { flags } = this.parse(this.constructor as any)
    // this.logLevel = flags.logLevel
  }
}
