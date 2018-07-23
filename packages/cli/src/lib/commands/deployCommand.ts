import { deployScenario } from '../deployScenario'
import * as inquirer from 'inquirer'
import * as fs from 'fs'
import * as ini from 'ini'
import * as Case from 'case'

import Locator from '../locationProvider'
import { ScenarioConfig } from '../types'

const deploy = (emitter, config: ScenarioConfig, locator: Locator) => async ({
  path = '.'
}) => {
  emitter.emit('deploy:started')
  const { BearerEnv } = config

  // Always true?
  if (!locator.scenarioRoot) {
    emitter.emit('rootPath:doesntExist')
    process.exit(1)
  }

  const mergedConfig = { ...config.bearerConfig, ...config.scenarioConfig }
  let { scenarioTitle } = mergedConfig
  const { OrgId } = mergedConfig

  const inquireScenarioTitle = () => {
    return inquirer.prompt([
      {
        message: 'Scenario title (e.g. attachPullRequest)?',
        type: 'input',
        name: 'scenarioTitle'
      }
    ])
  }
  if (!scenarioTitle) {
    emitter.emit('scenarioTitle:missing')
    try {
      const answers = await inquireScenarioTitle()
      scenarioTitle = answers.scenarioTitle
    } catch (e) {
      emitter.emit('scenarioTitle:creationFailed', e)
      process.exit(1)
    }
  }

  scenarioTitle = Case.kebab(Case.camel(scenarioTitle))
  const scenarioUuid = `${OrgId}-${scenarioTitle}`
  let scenarioConfigUpdate: any = { scenarioTitle }
  if (config.scenarioConfig.OrgId && config.scenarioConfig.OrgId !== OrgId) {
    scenarioConfigUpdate = {
      ...scenarioConfigUpdate,
      OrgId: config.scenarioConfig.OrgId
    }
  }

  fs.writeFileSync(locator.scenarioRc, ini.stringify(scenarioConfigUpdate))

  try {
    await deployScenario({ scenarioUuid }, emitter, config, locator)
    const setupUrl = `https://demo.bearer.tech/?scenarioUuid=${scenarioUuid}&scenarioTagName=${scenarioTitle}&name=${scenarioTitle}&orgId=${OrgId}&stage=${BearerEnv}`

    emitter.emit('deploy:finished', {
      scenarioUuid,
      scenarioTitle,
      setupUrl
    })
  } catch (error) {
    emitter.emit('deploy:failed', {
      error
    })
  }
}
module.exports = {
  useWith: (program, emitter, config, locator) => {
    program
      .command('deploy')
      .description(
        `Build a scenario package.
    $ bearer deploy
`
      )
      .action(deploy(emitter, config, locator))
  }
}
