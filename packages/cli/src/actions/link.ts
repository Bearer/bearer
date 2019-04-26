import * as inquirer from 'inquirer'

import BaseAction, { createExport } from './base'

class LinkAction extends BaseAction {
  async run(anIdentifier?: string) {
    let identifier = anIdentifier
    let { integrationTitle } = this.logger.bearerConfig as { integrationTitle: string }
    if (!identifier) {
      const { data } = await this.logger.devPortalClient.request<{ integrations: Integration[] }>({ query: QUERY })
      if (!data.data) {
        throw 'Unable to fetch integration list'
      }

      const { integrations } = data.data!
      const { integration } = await inquirer.prompt<{ integration: Integration }>([
        {
          name: 'integration',
          type: 'list',
          choices: integrations.map(int => ({
            name: `${int.uuid} - ${int.name}`,
            value: int
          }))
        }
      ])
      integrationTitle = integration.name
      identifier = integration.uuid
    }

    const integrationRc = { integrationTitle, integrationId: identifier }
    this.logger.bearerConfig.setIntegrationConfig(integrationRc)
    this.logger.log('Integration successfully linked! 🎉')
  }
}

export default createExport<[string] | []>(LinkAction)

type Integration = {
  uuid: string
  name: string
  latestActivity?: {
    state: string
  }
}

const QUERY = `
query CLILinkIntegrationList {
  integrations(includeGloballyAvailable: false) {
    uuid: uuidv2
    name
    latestActivity {
      state
    }
  }
}
`
