import BaseCommand from '../base-command'
import { RequireIntegrationFolder } from '../utils/decorators'

export default class Link extends BaseCommand {
  static description = 'Link your local integration to a remote one'

  static flags = {
    ...BaseCommand.flags
  }

  static args = [{ name: 'Integration_Identifier', required: true }]

  @RequireIntegrationFolder()
  async run() {
    const { args } = this.parse(Link)
    const identifier = args.Integration_Identifier
    const { integrationTitle } = this.bearerConfig
    const [orgId, integrationId] = identifier.replace(/\-/, '|').split('|')
    const integrationRc = { orgId, integrationId, integrationTitle }
    this.bearerConfig.setIntegrationConfig(integrationRc)
    this.log('Integration successfully linked! 🎉')
  }
}
