const deployScenario = require('../deployScenario')
const inquirer = require('inquirer')
const fs = require('fs')
const ini = require('ini')
const pathJs = require('path')
const Case = require('case')

const deploy = (emitter, config) => async ({ path = '.' }) => {
  emitter.emit('deploy:started')
  const { rootPathRc, BearerEnv } = config

  if (!rootPathRc) {
    emitter.emit('rootPath:doesntExist')
    process.exit(1)
  }

  const mergedConfig = { ...config.bearerConfig, ...config.scenarioConfig }
  let { scenarioTitle } = mergedConfig
  const { OrgId } = mergedConfig

  const inquireScenarioTitle = () => {
    emitter.emit('scenarioTitle:missing')
    return inquirer.prompt([
      {
        message: 'Scenario title (e.g. attachPullRequest)?',
        type: 'input',
        name: 'scenarioTitle'
      }
    ])
  }
  if (!scenarioTitle) {
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
  let scenarioConfigUpdate = { scenarioTitle }
  if (config.scenarioConfig.OrgId && config.scenarioConfig.OrgId !== OrgId) {
    scenarioConfigUpdate = {
      ...scenarioConfigUpdate,
      OrgId: config.scenarioConfig.OrgId
    }
  }

  fs.writeFileSync(pathJs.join(rootPathRc), ini.stringify(scenarioConfigUpdate))

  try {
    await deployScenario({ path, scenarioUuid }, emitter, config)
    const setupUrl = `https://demo.bearer.tech/?scenarioUuid=${scenarioUuid}&scenarioTagName=${scenarioTitle}&name=${scenarioTitle}&orgId=${OrgId}&stage=${BearerEnv}`

    emitter.emit('deploy:finished', {
      scenarioUuid,
      scenarioTitle,
      setupUrl
    })
  } catch (e) {
    console.log(e)
  }
}
module.exports = {
  useWith: (program, emitter, config) => {
    program
      .command('deploy')
      .description(
        `Build a scenario package.
    $ bearer deploy
`
      )
      .action(deploy(emitter, config))
  }
}
