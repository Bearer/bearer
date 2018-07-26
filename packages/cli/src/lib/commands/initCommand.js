const copy = require('copy-template-dir')
const path = require('path')
const Case = require('case')
const inquirer = require('inquirer')

const init = emitter => async scenarioTitle => {
  const { authenticationType } = await inquirer.prompt([
    {
      message: 'What kind of authentication do you want to use?',
      type: 'list',
      name: 'authenticationType',
      choices: [
        {
          name: 'OAuth2',
          value: 'oauth2'
        },
        {
          name: 'Basic Auth',
          value: 'basicauth'
        },
        {
          name: 'API Key',
          value: 'apikey'
        },
        {
          name: 'NoAuth',
          value: 'noauth'
        }
      ]
    }
  ])

  const vars = {
    scenarioTitle,
    componentName: Case.pascal(scenarioTitle),
    componentTagName: Case.kebab(scenarioTitle)
  }
  const inDir = path.join(__dirname, 'templates', 'init', authenticationType)
  const outDir = process.cwd()

  copy(inDir, outDir, vars, (err, createdFiles) => {
    if (err) throw err
    createdFiles.forEach(filePath => emitter.emit('generateTemplate:fileGenerated', filePath))
  })
}

module.exports = {
  useWith: (program, emitter, config) => {
    program
      .command('new <scenarioTitle>')
      .description(
        `Start a new scenario.
    $ bearer new myScenario
`
      )
      .action(init(emitter, config))
  }
}
