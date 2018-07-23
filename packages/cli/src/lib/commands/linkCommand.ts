import Locator from '../locationProvider'

export const link = (emitter, config, locator: Locator) => async scenarioId => {
  console.log(scenarioId)
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
