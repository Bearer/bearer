import BaseCommand from '../base-command'
import { getIntegrations, Integration } from '../utils/devPortal'

export async function linkIntegration(this: BaseCommand, anIdentifier: string) {
  let identifier = anIdentifier
  let { integrationTitle } = this.bearerConfig as { integrationTitle: string }
  if (!identifier) {
    const { integrations } = await getIntegrations(this)
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
