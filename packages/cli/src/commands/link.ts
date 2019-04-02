import BaseCommand from '../base-command'
import { RequireIntegrationFolder } from '../utils/decorators'
import { getIntegrations, Integration } from '../utils/devPortal'

export default class Link extends BaseCommand {
  static description = 'Link your local integration to a remote one'

  static flags = {
    ...BaseCommand.flags
  }

  static args = [{ name: 'Integration_Identifier' }]

  @RequireIntegrationFolder()
  async run() {
    const { args } = this.parse(Link)
    let identifier = args.Integration_Identifier
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
}
