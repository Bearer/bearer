import BaseCommand from '../base-command'
import { RequireIntegrationFolder } from '../utils/decorators'
import Link from '../actions/link'

export default class LinkCommand extends BaseCommand {
  static description = 'link to remote Bearer integration'

  static flags = {
    ...BaseCommand.flags
  }

  static args = [{ name: 'Integration_Identifier' }]

  @RequireIntegrationFolder()
  async run() {
    const { args } = this.parse(LinkCommand)
    await Link(this, args.Integration_Identifier)
  }
}
