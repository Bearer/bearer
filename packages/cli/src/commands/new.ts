import { flags } from '@oclif/command'
import * as Listr from 'listr'
import * as path from 'path'
import * as util from 'util'
const exec = util.promisify(require('child_process').exec)

import BaseCommand from '../BaseCommand'
import { AuthType } from '../types'

const authTypes = {
  [AuthType.OAuth2]: { name: 'OAuth2', value: AuthType.OAuth2 },
  [AuthType.Basic]: { name: 'Basic Auth', value: AuthType.Basic },
  [AuthType.ApiKey]: { name: 'API Key', value: AuthType.ApiKey },
  [AuthType.NoAuth]: { name: 'NoAuth', value: AuthType.NoAuth }
}

export default class New extends BaseCommand {
  static description = 'Generate a new scenario'

  static flags = {
    ...BaseCommand.flags,
    help: flags.help({ char: 'h' }),
    skipInstall: flags.boolean({ hidden: true }),
    authType: flags.string({
      char: 'a',
      description: 'Authorization type', // help description for flag
      options: Object.keys(authTypes).map(auth => authTypes[auth as AuthType].value) // only allow the value to be from a discrete set
    })
  }

  static args = [{ name: 'ScenarioName' }]

  async run() {
    const { args, flags } = this.parse(New)
    try {
      const name: string = args.ScenarioName || (await this.askForName())
      const authType: AuthType = (flags.authType as AuthType) || (await this.askForAuthType())
      const skipInstall = flags.skipInstall
      const tasks = new Listr([
        {
          title: 'Generating scenario files',
          task: (_ctx: any, _task: any) => {
            this.log('Generate files:')
            return this.copyFiles(name, authType)
          }
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
          task: async (_ctx: any, _task: any) => exec('yarn install', { cwd: path.join(this.copyDestFolder, name) })
        },
        {
          title: 'Installing scenario dependencies with npm',
          enabled: (ctx: any) => ctx.yarn === false && !skipInstall,
          task: () => exec('yarn install', { cwd: path.join(this.copyDestFolder, name) })
        }
      ])
      await tasks.run()
      this.success(`Scenario initialized, name: ${name}, authentication type: ${authTypes[authType].name}`)
      this.log("\nWhat's next?\n")
      this.log('* read the bearer documentation at https://docs.bearer.sh\n')
      this.log(`* start your local development environement by running:\n`)
      this.log(this.colors.bold(`   cd ${name} && bearer start`))
    } catch (e) {
      this.error('Could not generate new scenario' + e.toString())
    }
  }

  getVars = (name: string) => ({
    scenarioTitle: name,
    componentName: this.case.pascal(name),
    componentTagName: this.case.kebab(name),
    bearerTagVersion: process.env.BEARER_PACKAGE_VERSION || 'beta1'
  })

  get copyDestFolder(): string {
    return process.cwd()
  }

  async copyFiles(name: string, authType: AuthType) {
    await new Promise<boolean>((resolve, reject) => {
      const inDir = path.join(__dirname, '..', '..', 'templates', 'init', authType)
      this.debug(`Input directory: ${inDir}`)
      this.copy(inDir, this.copyDestFolder, this.getVars(name), (error, files) => {
        if (error) {
          this.error(error)
          reject(false)
        } else {
          this.printFiles(files)
          resolve(true)
        }
      })
    })
  }

  async askForName(): Promise<string> {
    const { name } = await this.inquirer.prompt<{ name: string }>([
      { message: 'Scenario name:', type: 'input', name: 'name' }
    ])
    return name.trim()
  }

  async askForAuthType(): Promise<AuthType> {
    const { authenticationType } = await this.inquirer.prompt<{ authenticationType: AuthType }>([
      {
        message: 'What kind of authentication do you want to use?',
        type: 'list',
        name: 'authenticationType',
        choices: Object.values(authTypes)
      }
    ])
    return authenticationType
  }

  printFiles(files: Array<string>) {
    const dest = this.copyDestFolder
    files.forEach(file => {
      this.log(this.colors.gray(`    create: `) + this.colors.white(file.replace(dest + '/', '')))
    })
    this.log('\n')
  }
}
