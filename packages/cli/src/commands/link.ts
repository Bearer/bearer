import BaseCommand from '../base-command'
import { RequireScenarioFolder } from '../utils/decorators'

export default class Link extends BaseCommand {
  static description = 'Link your local scenario to a remote one'

  static flags = {
    ...BaseCommand.flags
  }

  static args = [{ name: 'Scenario_Identifier', required: true }]

  @RequireScenarioFolder()
  async run() {
    const { args } = this.parse(Link)
    const identifier = args.Scenario_Identifier
    const { scenarioTitle } = this.bearerConfig
    const [orgId, scenarioId] = identifier.replace(/\-/, '|').split('|')
    const scenarioRc = { orgId, scenarioId, scenarioTitle }
    this.bearerConfig.setScenarioConfig(scenarioRc)
    this.log('Scenario successfully linked! 🎉')
  }
}
