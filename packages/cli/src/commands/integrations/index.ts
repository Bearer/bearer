import BaseCommand from '../../base-command'
import { ensureFreshToken } from '../../utils/decorators'

export default class IntegrationsList extends BaseCommand {
  static description = 'list deployed integrations'

  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  @ensureFreshToken()
  async run() {
    try {
      const { data } = await this.devPortalClient.request<{ integrations: Integration[] }>({
        query: QUERY
      })
      if (data.data) {
        const { integrations } = data.data
        const max = integrations.reduce(
          (acc, inte) => {
            acc.name = Math.max(inte.name.length, acc.name)
            acc.state = Math.max((inte.latestActivity || { state: '' }).state.length, acc.state)
            return acc
          },
          { name: 0, state: 0 }
        )
        this.log('')
        const headers = [
          {
            name: this.colors.bold('Name'.padEnd(max.name)),
            latestActivity: { state: this.colors.bold('Status'.padEnd(max.state)) }
          }
        ] as Integration[]
        headers.concat(integrations).forEach(inte => {
          this.log(
            '| %s | %s |',
            inte.name.padEnd(max.name),
            (inte.latestActivity || { state: '' }).state.padEnd(max.state)
          )
        })
      } else {
        this.error('Unable to fetch integration list')
      }
    } catch (e) {
      this.error(e)
    }
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

const QUERY = `
query CLIListIntegrations {
  integrations(includeGloballyAvailable: false) {
    deprecated_uuid: uuid
    uuid: uuidv2
    name
    latestActivity {
      state
    }
  }
}
`
