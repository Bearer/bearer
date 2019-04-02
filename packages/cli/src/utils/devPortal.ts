import axios from 'axios'
import BaseCommand from '../base-command'

export async function getIntegrations(command: BaseCommand) {
  const token = await command.bearerConfig.getToken()
  if (token) {
    const { data } = await axios.post<{ data: { integrations: Integration[] } }>(
      command.constants.DeveloperPortalAPIUrl,
      {
        query: QUERY
      },
      {
        headers: {
          Authorization: `Bearer ${token.id_token}`
        }
      }
    )
    return data.data
  }
  // TODO: manage error at the command level
  command.error('Unauthorized, please run bearer login')
  throw 'Error'
}

export type Integration = {
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
