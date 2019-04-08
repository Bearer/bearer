import { Authentications, TAuthentications } from '@bearer/types/lib/authentications'
import { flags } from '@oclif/command'
import * as fs from 'fs-extra'
import * as Listr from 'listr'
import * as path from 'path'

import BaseCommand from '../base-command'
import installDependencies from '../tasks/install-dependencies'
import initViews from '../tasks/init-views'
import { copyFiles, printFiles } from '../utils/helpers'

import GenerateComponent from './generate/component'
import GenerateSpec from './generate/spec'

export const authTypes: TAuthentications = {
  [Authentications.OAuth1]: { name: 'OAuth1', value: Authentications.OAuth1 },
  [Authentications.OAuth2]: { name: 'OAuth2', value: Authentications.OAuth2 },
  [Authentications.Basic]: { name: 'Basic Auth', value: Authentications.Basic },
  [Authentications.ApiKey]: { name: 'API Key', value: Authentications.ApiKey },
  [Authentications.NoAuth]: { name: 'NoAuth', value: Authentications.NoAuth },
  [Authentications.Custom]: { name: 'Custom', value: Authentications.Custom }
}

export default class New extends BaseCommand {
  static description = 'generate integration boilerplate'

  static flags = {
    ...BaseCommand.flags,
    help: flags.help({ char: 'h' }),
    skipInstall: flags.boolean({ hidden: true }),
    withViews: flags.boolean(),
    authType: flags.string({
      char: 'a',
      description: 'Authorization type', // help description for flag
      // only allow the value to be from a discrete set
      options: Object.keys(authTypes).map(auth => authTypes[auth as Authentications].value)
    })
  }

  static args = [{ name: 'IntegrationName' }]
  private destinationFolder!: string
  private path?: string

  async run() {
    const { args, flags } = this.parse(New)
    this.path = flags.path

    try {
      const name: string = args.IntegrationName || (await this.askForString('Integration name'))
      const authType: Authentications = (flags.authType as Authentications) || (await this.askForAuthType())
      const skipInstall = flags.skipInstall
      this.destinationFolder = name

      const tasks: Listr.ListrTask[] = [
        {
          title: 'Generating integration structure',
          task: async (ctx: any) => {
            try {
              const files = await this.createStructure(name, authType)
              ctx.files = files
              return true
            } catch (e) {
              this.error(e)
              return null
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
          title: 'Create views related files',
          enabled: () => flags.withViews,
          task: async (_ctx: any, _task: any) => {
            return new Listr([
              initViews({
                cmd: this,
                vars: this.initViewsVars(name, authType)
              }),
              {
                title: 'Create integration specification file',
                task: async (_ctx: any, _task: any) => GenerateSpec.run(['--path', this.copyDestFolder, '--silent'])
              },
              {
                title: 'Create intial components',
                task: async (_ctx: any, _task: any) =>
                  GenerateComponent.run(['feature', '--type', 'root', '--path', this.copyDestFolder, '--silent'])
              }
            ])
          }
        }
      ]

      if (!skipInstall) {
        tasks.push(installDependencies({ cwd: this.copyDestFolder }))
      }

      const { files } = await new Listr(tasks).run()
      printFiles(this, files)

      this.success(`Integration initialized, name: ${name}, authentication type: ${authTypes[authType].name}`)
      this.log("\nWhat's next?\n")
      this.log('* read the bearer documentation at https://docs.bearer.sh\n')
      this.log(`* start your local development environement by running:\n`)
      this.log(this.colors.bold(`   cd ${name} && bearer start`))
    } catch (e) {
      this.error(e)
    }
  }

  getVars = (name: string, authType: Authentications) => ({
    integrationTitle: name,
    componentName: this.case.pascal(name),
    componentTagName: this.case.kebab(name),
    bearerTagVersion: process.env.BEARER_PACKAGE_VERSION || 'beta6',
    bearerRestClient: this.bearerRestClient(authType)
  })

  bearerRestClient(authType: Authentications): string {
    switch (authType) {
      case Authentications.OAuth1:
        return '"oauth": "^0.9.15"'
      case Authentications.ApiKey:
      case Authentications.Basic:
      case Authentications.NoAuth:
      case Authentications.OAuth2:
      case Authentications.Custom:
        return '"axios": "^0.18.0"'
      default:
        throw new Error(`Authentication not found: ${authType}`)
    }
  }
  get copyDestFolder(): string {
    if (this.path) {
      return path.resolve(this.path)
    }
    return path.join(process.cwd(), this.destinationFolder)
  }

  createStructure(name: string, authType: Authentications): Promise<string[]> {
    if (fs.existsSync(this.copyDestFolder)) {
      return Promise.reject(this.colors.bold('Destination already exists: ') + this.copyDestFolder)
    }

    return copyFiles(this, path.join('init', 'structure'), this.copyDestFolder, this.getVars(name, authType), true)
  }

  initViewsVars = (name: string, authType: Authentications) => {
    const setup =
      authType === Authentications.NoAuth
        ? ''
        : `
    {
      classname: 'SetupAction',
      isRoot: true,
      initialTagName: 'setup-action',
      name: 'setup-action',
      label: 'Setup Action Component'
    },
    {
      classname: 'SetupView',
      isRoot: true,
      initialTagName: 'setup-view',
      name: 'setup-view',
      label: 'Setup Display Component'
    },`

    return { ...this.getVars(name, authType), setup }
  }

  createAuthenticationFiles(name: string, authType: Authentications): Promise<string[]> {
    return copyFiles(this, path.join('init', authType), this.copyDestFolder, this.getVars(name, authType), true)
  }

  async askForAuthType(): Promise<Authentications> {
    const { authenticationType } = await this.inquirer.prompt<{ authenticationType: Authentications }>([
      {
        message: 'Select an authentication method for the API you want to consume:',
        type: 'list',
        name: 'authenticationType',
        choices: Object.values(authTypes)
      }
    ])
    return authenticationType
  }
}
