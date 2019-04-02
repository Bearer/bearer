import { flags } from '@oclif/command'
import axios from 'axios'

import BaseCommand from '../../base-command'

export default class IntegrationsCreate extends BaseCommand {
  static flags = {
    ...BaseCommand.flags,
    description: flags.string({ char: 'd' }),
    name: flags.string({ char: 'n' })
  }

  async run() {
    const { flags } = this.parse(IntegrationsCreate)
    const name = flags.name || (await this.askForName())
    const description = flags.description || (await this.askForDescription())
    const token = await this.bearerConfig.getToken()

    if (token) {
      try {
        const response = await axios.post<{ data: { createIntegration: { integration: Integration } } }>(
          this.constants.DeveloperPortalAPIUrl,
          {
            query: MUTATION,
            variables: {
              name,
              description
            }
          },
          {
            headers: {
              Authorization: `Bearer ${token.id_token}`
            }
          }
        )
        const int = response.data.data.createIntegration.integration
        this.success('Integration successfully created')
        this.log(
          '      name: %s\n      uuid: %s\nidentifier: %s\n       Url:',
          int.name,
          int.uuid,
          int.deprecated_uuid,
          `${this.constants.DeveloperPortalUrl}integrations/${int.deprecated_uuid}`
        )
      } catch (e) {
        this.debug('%j', e.response)
        this.error(e)
      }
    }
  }

  async askForName(): Promise<string> {
    return this.askForString('Integration name')
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
