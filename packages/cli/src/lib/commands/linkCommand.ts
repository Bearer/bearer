import Locator from '../locationProvider'

export const link = (emitter, config, _locator: Locator) => async scenarioUuid => {
  emitter.emit('link:start')
  const { scenarioTitle } = config
  const [orgId, scenarioId] = scenarioUuid.replace(/\-/, '|').split('|')
  const scenarioRc = { orgId, scenarioId, scenarioTitle }
  config.setScenarioConfig(scenarioRc)
  emitter.emit('link:success', scenarioRc)
}
export function useWith(program, emitter, config, locator: Locator) {
  program
    .command('link')
    .description(
      `Link the scenario with developer portal
  $ bearer link 4l1c3-scenario-name
`
    )
    .action(link(emitter, config, locator))
}
