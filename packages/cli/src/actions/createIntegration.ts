import * as inquirer from 'inquirer'
// @ts-ignore
import * as suggest from 'inquirer-prompt-suggest'

import BaseAction, { createExport } from './base'
import { askForString } from '../utils/prompts'
import Link from './link'
import { randomNames } from '../utils/random'

class CreateIntegrationAction extends BaseAction {
  async run(options: TOptions = {}) {
    const name = options.name || (await this.askForName())
    const description = options.description || (await askForString('Description (optional)'))

    try {
      const { data } = await this.logger.devPortalClient.request<CreateIntegration>({
        query: MUTATION,
        variables: { name, description }
      })
      if (data.data) {
        const { integration } = data.data.createIntegration
        this.logger.success('Integration successfully created')
        this.logger.log(
          '      name: %s\n      uuid: %s\n       Url:',
          integration.name,
          integration.uuid,
          `${this.logger.constants.DeveloperPortalUrl}integrations/${integration.uuid}`
        )
        if (this.logger.isIntegrationLocation) {
          if (options.link) {
            await Link(this.logger, integration.uuid)
          }
        }
      } else {
        this.logger.debug(data)
        this.logger.error('Unable to create this integration, please retry')
      }
    } catch (e) {
      this.logger.debug('%j', e.response)
      this.logger.error(e)
    }
  }

  async askForName(): Promise<string> {
    inquirer.registerPrompt('suggest', suggest)
    const suggestions = randomNames(30)
    this.logger.debug('%j', suggestions)
    return askForString('Integration name', {
      suggestions,
      validate: (input: string) => {
        return input.length > 0
      },
      type: 'suggest'
    } as any)
  }
}

export default createExport<[TOptions]>(CreateIntegrationAction)

type Integration = {
  uuid: string
  name: string
  latestActivity?: {
    state: string
  }
}

type TOptions = { name?: string; description?: string; link?: boolean }

type CreateIntegration = { createIntegration: { integration: Integration } }

const MUTATION = `
mutation CLICreateIntegration($name: String!, $description: String!) {
  createIntegration(name: $name, description: $description) {
    integration {
      uuid: uuidv2
      name
    }
  }
}
`
