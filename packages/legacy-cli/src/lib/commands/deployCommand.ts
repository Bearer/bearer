import { deployScenario, IDeployOptions } from '../deployScenario'
import Locator from '../locationProvider'

const deploy = (emitter, config, locator: Locator) => async ({ viewsOnly = false, intentsOnly = false }) => {
  emitter.emit('deploy:started')
  // Always true?
  if (!locator.scenarioRoot) {
    emitter.emit('rootPath:doesntExist')
    process.exit(1)
  }

  const { hasScenarioLinked, scenarioId } = config
  if (!hasScenarioLinked) {
    emitter.emit('scenarioUuid:missing')
    process.exit(1)
  }

  if (
    !(
      config.bearerConfig &&
      config.bearerConfig.authorization &&
      config.bearerConfig.authorization.AuthenticationResult &&
      config.bearerConfig.authorization.AuthenticationResult.IdToken
    )
  ) {
    emitter.emit('user:notAuthenticated')
    process.exit(1)
  }

  const deployOptions: IDeployOptions = { noViews: intentsOnly, noIntents: viewsOnly }

  try {
    await deployScenario(deployOptions, emitter, config, locator)
    const setupUrl = `${config.DeveloperPortalUrl}scenarios/${config.scenarioUuid}/preview`

    emitter.emit('deploy:finished', {
      scenarioId,
      setupUrl
    })
  } catch (error) {
    emitter.emit('deploy:failed', {
      error
    })
  }
}

export function useWith(program, emitter, config, locator): void {
  program
    .command('deploy')
    .description(
      `Build a scenario package.
$ bearer deploy
`
    )
    .option('-v, --views-only', 'Deploy views only')
    .option('-i, --intents-only', 'Deploy intents only')
    .action(deploy(emitter, config, locator))
}
