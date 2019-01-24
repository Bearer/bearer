import { flags } from '@oclif/command'
import * as fs from 'fs'
import * as path from 'path'

import BaseCommand from '../../base-command'
import { RequireScenarioFolder } from '../../utils/decorators'
import { copyFiles } from '../../utils/helpers'
import * as Listr from 'listr'

export default class GenerateSetup extends BaseCommand {
  static description = 'Generate a Bearer Setup'
  static hidden = true
  static flags = { ...BaseCommand.flags, force: flags.boolean({}) }

  static args = []

  @RequireScenarioFolder()
  async run() {
    const { flags } = this.parse(GenerateSetup)
    if (flags.force || !setupExists(this.locator.srcViewsDir)) {
      const fields = this.scenarioAuthConfig.setupViews
      if (fields && fields.length) {
        try {
          const tasks: Listr.ListrTask[] = [
            {
              title: 'Generating setup components',
              task: async () => {
                try {
                  await copyFiles(
                    this,
                    'generate/setup',
                    this.locator.srcViewsDir,
                    this.getVars(this.bearerConfig.scenarioConfig.scenarioTitle, fields, 'NONE')
                  )
                  return true
                } catch (e) {
                  this.error(e)
                  return null
                }
              }
            },
            {
              title: 'Generating setup intents',
              task: async () => {
                try {
                  console.log(this.scenarioAuthConfig.authType)
                  await copyFiles(
                    this,
                    `generate/setup-intents`,
                    this.locator.srcIntentsDir,
                    this.getVars(
                      this.bearerConfig.scenarioConfig.scenarioTitle,
                      fields,
                      this.scenarioAuthConfig.authType
                    )
                  )
                  return true
                } catch (e) {
                  this.error(e)
                  return null
                }
              }
            }
          ]
          await new Listr(tasks).run()
        } catch (e) {
          this.error(e)
        }
      } else {
        this.warn(`No setupViews key found within auth.config.json file. skipping`)
      }
    } else {
      this.warn('Setup files already exist, use --force flag to overwrite existing ones.')
    }
  }

  getVars(scenarioTitle: string, fields: any, authType: string) {
    const contextAuthType = getContextAuthType(authType)
    return {
      contextAuthType,
      componentName: this.case.pascal(scenarioTitle),
      componentTagName: this.case.kebab(scenarioTitle),
      fields: JSON.stringify(fields),
      contextAuthTypeImport: contextAuthType ? `${contextAuthType}, ` : '',
      contextAuthTypeImplements: contextAuthType ? `, ${contextAuthType}` : '',
      dataType: getDataType(authType)
    }
  }
}

function getDataType(authType: string): string {
  switch (authType) {
    case 'NONE':
      return '{}'
    case 'BASIC':
      return '{username: string, password: string}'
    case 'APIKEY':
      return '{apiKey: string}'
    case 'OAUTH2':
      return '{accessToken: string}'
    default:
      return 'any'
  }
}

function getContextAuthType(authType: string): string {
  switch (authType) {
    case 'NONE':
      return ''
    case 'BASIC':
      return 'TBASICAuthContext'
    case 'APIKEY':
      return 'TAPIKEYAuthContext'
    case 'OAUTH2':
      return 'TOAUTH2AuthContext'
    default:
      return 'any'
  }
}
// Note: using Or condition in case the developer delete one but customized the other component
function setupExists(location: string): boolean {
  return (
    fs.existsSync(path.join(location, 'setup-action.tsx')) || fs.existsSync(path.join(location, 'setup-display.tsx'))
  )
}
