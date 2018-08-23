import Authentications from '@bearer/types/lib/Authentications'
import { flags } from '@oclif/command'
import * as fs from 'fs-extra'
import * as Listr from 'listr'
import * as path from 'path'
import * as util from 'util'

import BaseCommand from '../BaseCommand'
import { copyFiles, printFiles } from '../utils/helpers'

import GenerateComponent from './generate/component'
import GenerateSetup from './generate/setup'
import GenerateSpec from './generate/spec'

const authTypes = {
  [Authentications.OAuth2]: { name: 'OAuth2', value: Authentications.OAuth2 },
  [Authentications.Basic]: { name: 'Basic Auth', value: Authentications.Basic },
  [Authentications.ApiKey]: { name: 'API Key', value: Authentications.ApiKey },
  [Authentications.NoAuth]: { name: 'NoAuth', value: Authentications.NoAuth }
}

const exec = util.promisify(require('child_process').exec)

export default class New extends BaseCommand {
  static description = 'Generate a new scenario'

  static flags = {
    ...BaseCommand.flags,
    help: flags.help({ char: 'h' }),
    skipInstall: flags.boolean({ hidden: true }),
    authType: flags.string({
      char: 'a',
      description: 'Authorization type', // help description for flag
      options: Object.keys(authTypes).map(auth => authTypes[auth as Authentications].value) // only allow the value to be from a discrete set
    })
  }

  static args = [{ name: 'ScenarioName' }]
  private destinationFolder!: string
  private path?: string

  async run() {
    const { args, flags } = this.parse(New)
    this.path = flags.path

    try {
      const name: string = args.ScenarioName || (await this.askForString('Scenario name:'))
      const authType: Authentications = (flags.authType as Authentications) || (await this.askForAuthType())
      const skipInstall = flags.skipInstall

      const tasks = new Listr([
        {
          title: 'Generating scenario structure',
          task: async (ctx: any) => {
            try {
              const files = await this.createStructure(name)
              ctx.files = files
              return true
            } catch (e) {
              this.error(e)
            }
          }
        },
        {
          title: 'Create auth files',
          task: async (ctx: any, _task: any) => {
            const files = await this.createAuthenticationFiles(name, authType)
            ctx.files = [...ctx.files, ...files]
          }
        },
        {
          title: 'Create setup files',
          task: (_ctx: any, _task: any) => GenerateSetup.run(['--path', this.copyDestFolder, '--silent'])
        },
        {
          title: 'Create scenario specification file',
          task: (_ctx: any, _task: any) => GenerateSpec.run(['--path', this.copyDestFolder, '--silent'])
        },
        {
          title: 'Create intial components',
          task: (_ctx: any, _task: any) =>
            GenerateComponent.run(['feature', '--type', 'root', '--path', this.copyDestFolder, '--silent'])
        },
        {
          title: 'npm/yarn lookup',
          enabled: (_ctx: any) => !skipInstall,
          task: async (ctx: any, task: any) =>
            exec('yarn -v')
              .then(() => {
                ctx.yarn = true
              })
              .catch(() => {
                ctx.yarn = false
                task.skip('Yarn not available, install it via `npm install -g yarn`')
              })
        },
        {
          title: 'Installing scenario dependencies with yarn',
          enabled: (ctx: any) => ctx.yarn === true && !skipInstall,
          task: async (_ctx: any, _task: any) => exec('yarn install', { cwd: this.copyDestFolder })
        },
        {
          title: 'Installing scenario dependencies with npm',
          enabled: (ctx: any) => ctx.yarn === false && !skipInstall,
          task: async () => exec('yarn install', { cwd: this.copyDestFolder })
        }
      ])

      const { files } = await tasks.run()
      printFiles(this, files)

      this.success(`Scenario initialized, name: ${name}, authentication type: ${authTypes[authType].name}`)
      this.log("\nWhat's next?\n")
      this.log('* read the bearer documentation at https://docs.bearer.sh\n')
      this.log(`* start your local development environement by running:\n`)
      this.log(this.colors.bold(`   cd ${name} && bearer start`))
    } catch (e) {
      this.error(e)
    }
  }

  getVars = (name: string) => ({
    scenarioTitle: name,
    componentName: this.case.pascal(name),
    componentTagName: this.case.kebab(name),
    bearerTagVersion: process.env.BEARER_PACKAGE_VERSION || 'beta1'
  })

  get copyDestFolder(): string {
    if (this.path) {
      return path.resolve(this.path)
    }
    return path.join(process.cwd(), this.destinationFolder)
  }

  createStructure(name: string): Promise<Array<string>> {
    this.destinationFolder = name
    if (fs.existsSync(this.copyDestFolder)) {
      return Promise.reject(this.colors.bold('Destination already exists: ') + this.copyDestFolder)
    }
    return copyFiles(this, path.join('init', 'structure'), this.copyDestFolder, this.getVars(name), true)
  }

  createAuthenticationFiles(name: string, authType: Authentications): Promise<Array<string>> {
    return copyFiles(this, path.join('init', authType), this.copyDestFolder, this.getVars(name), true)
  }

  async askForAuthType(): Promise<Authentications> {
    const { authenticationType } = await this.inquirer.prompt<{ authenticationType: Authentications }>([
      {
        message: 'What kind of authentication do you want to use?',
        type: 'list',
        name: 'authenticationType',
        choices: Object.values(authTypes)
      }
    ])
    return authenticationType
  }
}
