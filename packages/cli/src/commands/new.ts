import { Authentications, TAuthentications } from '@bearer/types/lib/authentications'
import { flags } from '@oclif/command'
import * as fs from 'fs-extra'
import * as Listr from 'listr'
import * as path from 'path'
import * as Case from 'case'
import * as os from 'os'
import * as globby from 'globby'
import cliUx from 'cli-ux'
import { exec } from 'child_process'
import * as util from 'util'

import BaseCommand from '../base-command'
import installDependencies from '../tasks/install-dependencies'
import initViews from '../tasks/init-views'
import { copyFiles, printFiles } from '../utils/helpers'

import GenerateComponent from './generate/component'
import GenerateSpec from './generate/spec'
import * as inquirer from 'inquirer'
import { askForString } from '../utils/prompts'
import { INTEGRATION_PROOF } from '../utils/locator'

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
    template: flags.string({
      char: 't',
      description: 'Generate an integration from a template (git url)'
    }),
    directory: flags.string({
      char: 'd',
      dependsOn: ['template'],
      description: 'Select a directory as source of the integration'
    }),
    skipInstall: flags.boolean({ hidden: true, description: 'Do not install dependencies' }),
    withViews: flags.boolean({ description: 'Experimental - generate views' }),
    force: flags.boolean({ description: 'Force copying files', char: 'f' }),
    authType: flags.string({
      char: 'a',
      description: 'Authorization type', // help description for flag
      // only allow the value to be from a discrete set
      options: Object.keys(authTypes).map(auth => authTypes[auth as Authentications].value)
    })
  }

  static args = [{ name: 'IntegrationName' }]
  copyDestinationFolder!: string
  private path?: string
  private name!: string

  async run() {
    const { args, flags } = this.parse(New)
    this.debug('url: %j', flags)
    const { template, directory } = flags
    this.path = flags.path

    try {
      this.name = args.IntegrationName || (await askForString('Integration name'))
      this.copyDestinationFolder = await defineLocationPath(this, {
        name: this.name,
        cwd: this.path || process.cwd(),
        force: flags.force
      })

      this.debug('target path: %s', this.copyDestinationFolder)
      const skipInstall = flags.skipInstall
      let files: string[] = []

      if (template) {
        const tmp = path.join(os.tmpdir(), Date.now().toString())
        await cloneRepository(template, tmp, this)

        const { selected } = await selectFolder(tmp, { selectedPath: directory })

        const source = path.join(tmp, selected)
        fs.copySync(source, this.copyDestinationFolder)
        this.debug('copied from %s to %s', source, this.copyDestinationFolder)
        files = await globby(['**/*'], { cwd: this.copyDestinationFolder })
      } else {
        const authType: Authentications = (flags.authType as Authentications) || (await askForAuthType())
        const tasks: Listr.ListrTask[] = [
          {
            title: 'Generating integration structure',
            task: this.createIntegrationStructure(authType)
          },
          {
            title: 'Create auth files',
            task: this.createAuthFiles(authType)
          },
          {
            title: 'Create views related files',
            enabled: () => flags.withViews,
            task: this.createViewsStructure(authType)
          }
        ]

        const { files: created } = await new Listr(tasks, { nonTTYRenderer: 'silent' }).run()
        files = created
        if (authType) {
          this.success(`Integration initialized, name: ${this.name}, authentication type: ${authTypes[authType].name}`)
        }
      }

      if (!skipInstall) {
        await new Listr([installDependencies({ cwd: this.copyDestinationFolder })]).run()
      }

      printFiles(this, files)
      const finalLocation = this.copyDestinationFolder.replace(path.resolve(args.path || process.cwd()), '.')
      this.log("\nWhat's next?\n")
      this.log('* read the bearer documentation at https://docs.bearer.sh\n')
      this.log(`* start your local development environment by running:\n`)
      this.log(this.colors.bold(`   cd ${finalLocation}`))
    } catch (e) {
      this.debug('error: %j', e)
      if (e instanceof NoIntegrationFoundError) {
        this.error(this.colors.red('No valid integrations found within the cloned git repository'))
      } else {
        this.error(e)
      }
    }
  }

  private createIntegrationStructure = (authType: Authentications) => async (ctx: { files: string[] }) => {
    ctx.files = await copyFiles(
      this,
      path.join('init', 'structure'),
      this.copyDestinationFolder,
      this.getVars(authType),
      true
    )
  }

  private createAuthFiles = (authType: Authentications) => async (ctx: { files: string[] }) => {
    const files = await copyFiles(
      this,
      path.join('init', authType),
      this.copyDestinationFolder,
      this.getVars(authType),
      true
    )
    ctx.files = [...ctx.files, ...files]
  }

  private createViewsStructure = (authType: Authentications) => async (_ctx: any, _task: any) => {
    return new Listr([
      initViews({
        cmd: this,
        vars: { ...this.getVars(authType), ...initViewsVars(authType) }
      }),
      {
        title: 'Create integration specification file',
        task: async (_ctx: any, _task: any) => GenerateSpec.run(['--path', this.copyDestinationFolder, '--silent'])
      },
      {
        title: 'Create intial components',
        task: async (_ctx: any, _task: any) =>
          GenerateComponent.run(['feature', '--type', 'root', '--path', this.copyDestinationFolder, '--silent'])
      }
    ])
  }

  private getVars = (authType: Authentications) => ({
    integrationTitle: this.name,
    componentName: Case.pascal(this.name),
    componentTagName: Case.kebab(this.name),
    bearerTagVersion: process.env.BEARER_PACKAGE_VERSION || 'latest',
    bearerRestClient: bearerRestClient(authType)
  })
}

export async function askForAuthType(): Promise<Authentications> {
  const { authenticationType } = await inquirer.prompt<{ authenticationType: Authentications }>([
    {
      message: 'Select an authentication method for the API you want to consume:',
      type: 'list',
      name: 'authenticationType',
      choices: Object.values(authTypes)
    }
  ])
  return authenticationType
}

export function bearerRestClient(authType: Authentications): string {
  switch (authType) {
    case Authentications.OAuth1:
      return '"oauth": "^0.9.15"'
    case Authentications.ApiKey:
    case Authentications.Basic:
    case Authentications.NoAuth:
    case Authentications.OAuth2:
    case Authentications.Custom:
      return '"axios": "^0.18.0"'
  }
  throw new Error(`Authentication not found: ${authType}`)
}
// @ts-ignore
export async function cloneRepository(url: string, destination: string, logger: BaseCommand) {
  const asyncExec = util.promisify(exec)
  cliUx.action.start('cloning')
  try {
    const { stdout, stderr } = await asyncExec('git --version')
    // @ts-ignore
    logger.debug('out: %j\n, err: %j', stdout, stderr)
  } catch (e) {
    // @ts-ignore
    logger.debug('error: %j', e)
    cliUx.action.stop()
    throw new GitNotAvailable()
  }

  await new Promise(async (resolve, reject) => {
    try {
      const command = `git clone ${url} --depth=1 ${destination}`
      // @ts-ignore
      logger.debug(`Running ${command}`)
      await asyncExec(command)
      resolve()
    } catch (e) {
      cliUx.action.stop()
      // @ts-ignore
      logger.debug(`%j`, e)
      reject(new CloningError())
    }
  })
  cliUx.action.stop()
}

export async function selectFolder(
  location: string,
  { selectedPath, integrationRootProof = INTEGRATION_PROOF }: { integrationRootProof?: string; selectedPath?: string }
) {
  const list = await globby([`**/${integrationRootProof}`], { cwd: location })
  const choices = list
    .sort()
    .map(path => {
      return {
        name: path.split(`/${integrationRootProof}`)[0],
        value: path.split(integrationRootProof)[0].replace(/\/$/, '')
      }
    })
    .filter(choice => {
      return !Boolean(selectedPath) || choice.value === selectedPath
    })

  switch (choices.length) {
    case 0: {
      throw new NoIntegrationFoundError(selectedPath)
    }
    case 1: {
      return { selected: choices[0].value }
    }
    default: {
      return await inquirer.prompt<{ selected: string }>([
        {
          choices,
          name: 'selected',
          message: 'Select a template',
          type: 'list'
        }
      ])
    }
  }
}

export async function defineLocationPath(
  logger: { warn: (message: string) => void },
  { name, cwd, force }: { name: string; cwd: string; force?: boolean | undefined }
) {
  let location = path.join(cwd, name)
  let shouldForce = force
  while (fs.existsSync(location) && shouldForce !== true) {
    if (shouldForce === false) {
      logger.warn(`${location.replace(cwd, '')} already exists\n please provide a different folder name`)
      const folderName = await askForString('Folder name to use')
      location = path.join(cwd, folderName)
    } else {
      const { override } = await inquirer.prompt<{ override: boolean }>([
        {
          type: 'confirm',
          name: 'override',
          default: false,
          message: `${location} is not empty, copy files anyway?`
        }
      ])
      shouldForce = override
    }
  }
  return location
}

export function initViewsVars(authType: Authentications) {
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

  return { setup }
}

/**
 * Custom errors
 */

class NoIntegrationFoundError extends Error {
  constructor(selectedFolder?: string) {
    super(`No valid integrations found ${selectedFolder ? `under ${selectedFolder}` : ''}`)
  }
}

class GitNotAvailable extends Error {
  constructor() {
    super('git command not found in your path, please install it')
  }
}

class CloningError extends Error {
  constructor() {
    super('Error while cloning the repository')
  }
}
