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
    identifier = integration.deprecated_uuid
  }

  const [orgId, integrationId] = identifier.replace(/\-/, '|').split('|')
  const integrationRc = { orgId, integrationId, integrationTitle }
  this.bearerConfig.setIntegrationConfig(integrationRc)
  this.log('Integration successfully linked! 🎉')
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
query CLILinkIntegrationList {
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
