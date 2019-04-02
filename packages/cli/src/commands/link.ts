import BaseCommand from '../base-command'
import { RequireIntegrationFolder } from '../utils/decorators'
import { linkIntegration } from '../utils/commands'

export default class Link extends BaseCommand {
  static description = 'Link your local integration to a remote one'

  static flags = {
    ...BaseCommand.flags
  }

  static args = [{ name: 'Integration_Identifier' }]

  @RequireIntegrationFolder()
  async run() {
    const { args } = this.parse(Link)
    await linkIntegration.bind(this)(args.Integration_Identifier)
  }
}
