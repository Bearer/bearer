import Authentications from '@bearer/types/lib/authentications'
import { flags } from '@oclif/command'
import * as Case from 'case'

import BaseCommand from '../../base-command'
import { RequireIntegrationFolder, skipIfNoViews } from '../../utils/decorators'
import { copyFiles } from '../../utils/helpers'
import { askForString } from '../../utils/prompts'
import * as inquirer from 'inquirer'

enum TComponent {
  BLANK = 'blank',
  COLLECTION = 'collection',
  ROOT = 'root'
}

export default class GenerateComponent extends BaseCommand {
  static description = 'Generate a Bearer component'
  static aliases = ['g:c']
  static flags = {
    ...BaseCommand.flags,
    type: flags.string({ char: 't', options: Object.values(TComponent) })
  }

  static args = [{ name: 'name' }]

  @skipIfNoViews(true)
  @RequireIntegrationFolder()
  async run() {
    const { args, flags } = this.parse(GenerateComponent)
    const type: TComponent = (flags.type as TComponent) || (await this.askForComponentType())
    const name: string = args.name || (await askForString('Name'))
    const outDir = type === TComponent.ROOT ? this.locator.srcViewsDir : this.locator.srcViewsDirResource('components')

    try {
      await copyFiles(
        this,
        `generate/${type}Component`,
        outDir,
        this.getVars(name, this.integrationAuthConfig.authType)
      )
      // TODO: add a nicer display
      this.success(`\nComponent generated`)
      if (type === TComponent.ROOT) {
        this.warn(
          this.colors.italic(
            // tslint:disable-next-line:max-line-length
            'Please make sure to update the spec.ts file to reflect your new Root Component on the Developer Portal preview page.'
          )
        )
      }
    } catch (e) {
      this.error(e)
    }
  }

  getVars(name: string, authType: Authentications) {
    const componentName = Case.pascal(name)
    return {
      componentName,
      fileName: name,
      componentClassName: componentName, // it gives more meaning within templates
      componentTagName: Case.kebab(componentName),
      groupName: Case.kebab(componentName),
      withAuthScreen: authType === Authentications.OAuth2 ? '<bearer-navigator-auth-screen />' : null
    }
  }

  async askForComponentType(): Promise<TComponent> {
    const { type } = await inquirer.prompt<{ type: TComponent }>([
      {
        choices,
        message: 'What kind of component would you like to generate:',
        type: 'list',
        name: 'type'
      }
    ])
    return type
  }
}

// TODO: better names
const choices = [
  {
    name: 'Blank',
    value: TComponent.BLANK
  },
  {
    name: 'Collection',
    value: TComponent.COLLECTION
  },
  {
    name: 'Root component',
    value: TComponent.ROOT
  }
]
