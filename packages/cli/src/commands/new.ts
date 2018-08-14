import { flags } from '@oclif/command'
import * as path from 'path'

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
    authType: flags.string({
      char: 'a',
      description: 'Authorization type', // help description for flag
      options: Object.keys(authTypes).map(auth => authTypes[auth as AuthType].value) // only allow the value to be from a discrete set
    })
  }

  static args = [{ name: 'ScenarioName' }]

  async run() {
    this.ux.action.start('Generating scenario files')
    const { args, flags } = this.parse(New)
    try {
      const name: string = args.ScenarioName || (await this.askForName())
      const authType: AuthType = (flags.authType as AuthType) || (await this.askForAuthType())
      await this.copyFiles(name, authType)
      await this.ux.wait()
      this.ux.action.stop()
      this.log(`Scenario initialized, name: ${name}, authentication type: ${authTypes[authType].name}`)
    } catch (e) {
      this.error('Could not generate new scenario' + e.toString())
    }
  }

  getVars = (name: string) => ({
    scenarioTitle: name,
    componentName: this.case.pascal(name),
    componentTagName: this.case.kebab(name)
  })

  async copyFiles(name: string, authType: AuthType) {
    this.log(`Copying initial files for ${authType}`)
    await new Promise<boolean>((resolve, reject) => {
      const inDir = path.join(__dirname, '..', '..', 'templates', 'init', authType)
      this.debug(`Input directory: ${inDir}`)
      this.copy(inDir, process.cwd(), this.getVars(name), (error, files) => {
        if (error) {
          this.error(error)
          reject(false)
        } else {
          this.debug(`Successfully copied files:\n\t* ${files.join('\n\t* ')}`)
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
}
