import BaseCommand from '../../base-command'
import axios from 'axios'

export default class IntegrationsList extends BaseCommand {
  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  async run() {
    const token = await this.bearerConfig.getToken()
    if (token) {
      try {
        const response = await axios.post<{ data: { integrations: Integration[] } }>(
          this.constants.DeveloperPortalAPIUrl,
          {
            query: QUERY
          },
          {
            headers: {
              Authorization: `Bearer ${token.id_token}`
            }
          }
        )
        const max = response.data.data.integrations.reduce(
          (acc, inte) => {
            acc.name = Math.max(inte.name.length, acc.name)
            acc.name = Math.max((inte.latestActivity || { state: '' }).state.length, acc.name)
            acc.uuid = Math.max(inte.uuid.length, acc.uuid)
            return acc
          },
          { name: 0, state: 0, uuid: 0 }
        )

        response.data.data.integrations.forEach(inte => {
          this.log(
            '| %s | %s | %s',
            inte.name.padEnd(max.name),
            (inte.latestActivity || { state: '' }).state.padEnd(max.state),
            inte.uuid.padEnd(max.uuid)
          )
        })
      } catch (e) {
        console.log(e.request)
      }
    }
  }
}

type Integration = {
  uuid: string
  name: string
  latestActivity?: {
    state: string
  }
}

const QUERY = `
query CLIListIntegrations {
  integrations(includeGloballyAvailable: false) {
    identifier
    uuid
    name
    latestActivity {
      state
    }
  }
}
`
