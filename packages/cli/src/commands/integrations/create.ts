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
        const response = await axios.post<{ data: { integrations: Integration[] } }>(
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
        this.success(JSON.stringify(response.data))
      } catch (e) {
        console.log(e)
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
      uuidv2
      uuid
      name
    }
  }
} 
`
