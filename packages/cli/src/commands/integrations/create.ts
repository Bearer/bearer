import { flags } from '@oclif/command'

import BaseCommand from '../../base-command'
import { linkIntegration } from '../../utils/commands'
import { ensureFreshToken } from '../../utils/decorators'

export default class IntegrationsCreate extends BaseCommand {
  static flags = {
    ...BaseCommand.flags,
    description: flags.string({ char: 'd' }),
    name: flags.string({ char: 'n' }),
    skipLink: flags.boolean({ char: 'l' })
  }

  @ensureFreshToken()
  async run() {
    const { flags } = this.parse(IntegrationsCreate)
    const name = flags.name || (await this.askForName())
    const description = flags.description || (await this.askForDescription())

    try {
      const { data } = await this.devPortalClient.request<CreateIntegration>({
        query: MUTATION,
        variables: { name, description }
      })
      if (data.data) {
        const { integration } = data.data.createIntegration
        this.success('Integration successfully created')
        this.log(
          '      name: %s\n      uuid: %s\nidentifier: %s\n       Url:',
          integration.name,
          integration.uuid,
          integration.deprecated_uuid,
          `${this.constants.DeveloperPortalUrl}integrations/${integration.deprecated_uuid}`
        )
        if (this.bearerConfig.isIntegrationLocation) {
          // tslint:disable-next-line no-boolean-literal-compare
          if (!flags.skipLink) {
            await linkIntegration.bind(this)(integration.deprecated_uuid)
          }
        }
      } else {
        this.debug(data)
        this.error('Unable to create this integration, please retry')
      }
    } catch (e) {
      this.debug('%j', e.response)
      this.error(e)
    }
  }

  async askForName(): Promise<string> {
    return this.askForString('Integration name', {
      validate: (input: string) => {
        return input.length > 0
      }
    })
  }
  async askForDescription(): Promise<string> {
    return this.askForString('Description (optional)')
  }
}

type Integration = {
  deprecated_uuid: string
  uuid: string
  name: string
  latestActivity?: {
    state: string
  }
}

type CreateIntegration = { createIntegration: { integration: Integration } }

const MUTATION = `
mutation CLICreateIntegration($name: String!, $description: String!) {
  createIntegration(name: $name, description: $description) {
    integration {
      deprecated_uuid: uuid
      uuid: uuidv2
      name
    }
  }
} 
`
