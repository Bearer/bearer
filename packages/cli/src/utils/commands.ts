import BaseCommand from '../base-command'

export async function linkIntegration(this: BaseCommand, anIdentifier: string) {
  let identifier = anIdentifier
  let { integrationTitle } = this.bearerConfig as { integrationTitle: string }
  if (!identifier) {
    const { data } = await this.devPortalClient.request<{ integrations: Integration[] }>({ query: QUERY })
    if (!data.data) {
      throw 'Unable to fetch integration list'
    }

    const { integrations } = data.data!
    const { integration } = await this.inquirer.prompt<{ integration: Integration }>([
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
  this.bearerConfig.setIntegrationConfig(integrationRc)
  this.log('Integration successfully linked! 🎉')
}

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
