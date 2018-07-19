const copy = require('copy-template-dir')
const del = require('del')
const path = require('path')
const inquirer = require('inquirer')
const Case = require('case')
const intents = require('@bearer/intents')
const templates = require('@bearer/templates')
const rc = require('rc')

const INTENT = 'intent'
const SCREEN = 'screen'

async function generateTemplates({ emitter, templateType, rootPathRc }) {
  const authConfig = require(path.join(
    path.dirname(rootPathRc),
    'intents',
    'auth.config.json'
  ))

  const scenarioConfig = rc('scenario')
  const { scenarioTitle } = scenarioConfig

  const configKey = `${templateType}Screens`

  const vars = {
    scenarioTitle: Case.camel(scenarioTitle),
    componentTagName: Case.kebab(scenarioTitle),
    fields: authConfig[configKey] ? JSON.stringify(authConfig[configKey]) : '[]'
  }
  const inDir = path.join(__dirname, `templates/generate/${templateType}`)
  const outDir = path.join(path.dirname(rootPathRc), '/screens/.build/src/')

  await del(`${outDir}*${templateType}*.tsx`).then(paths => {
    console.log('Deleted files and folders:\n', paths.join('\n'))
  })

  copy(inDir, outDir, vars, (err, createdFiles) => {
    if (err) throw err
    createdFiles.forEach(filePath =>
      emitter.emit('generateIntent:fileGenerated', filePath)
    )
  })
}

const generate = (emitter, { rootPathRc }) => async env => {
  if (!rootPathRc) {
    emitter.emit('rootPath:doesntExist')
    process.exit(1)
  }

  if (env.config) {
    return generateTemplates({ emitter, templateType: 'config', rootPathRc })
  }

  if (env.setup) {
    return generateTemplates({ emitter, templateType: 'setup', rootPathRc })
  }

  const { template } = await inquirer.prompt([
    {
      message: 'What do you want to generate',
      type: 'list',
      name: 'template',
      choices: [
        {
          name: 'Intent',
          value: INTENT
        },
        {
          name: 'Screen',
          value: SCREEN
        }
      ]
    }
  ])

  const params = { emitter, rootPathRc }

  switch (template) {
    case INTENT:
      generateIntent(params)
      break
    case SCREEN:
      await generateScreen(params)
      break
    default:
  }
}

async function askForName() {
  const { name } = await inquirer.prompt([
    {
      message: 'Give it a name',
      type: 'input',
      name: 'name'
    }
  ])

  return name
}

async function generateScreen({ emitter, rootPathRc }) {
  const name = await askForName()
  const componentName = Case.pascal(name)
  const vars = {
    screenName: componentName,
    componentTagName: Case.kebab(componentName)
  }
  const inDir = path.join(__dirname, 'templates/generate/screen')
  const outDir = path.join(path.dirname(rootPathRc), 'screens/src/components')

  copy(inDir, outDir, vars, (err, createdFiles) => {
    if (err) throw err
    createdFiles.forEach(filePath =>
      emitter.emit('generateIntent:fileGenerated', filePath)
    )
  })
}

const filteredChoices = (intents, propertyFlag) =>
  Object.keys(intents)
    .filter(intent => intents[intent][propertyFlag])
    .map(intent => ({
      name: intents[intent].display,
      value: intent
    }))
    .sort((a, b) => a.name > b.name)

const choices = [
  ...filteredChoices(intents, 'isGlobalIntent'),
  new inquirer.Separator(),
  ...filteredChoices(intents, 'isStateIntent')
]

function getActionExample(intentType, authType) {
  return templates[authType][intentType]
}

async function generateIntent({ emitter, rootPathRc }) {
  const { intentType } = await inquirer.prompt([
    {
      message: 'What type of intent do you want to generate',
      type: 'list',
      name: 'intentType',
      choices
    }
  ])
  const name = await askForName()
  const authConfig = require(path.join(
    path.dirname(rootPathRc),
    'intents',
    'auth.config.json'
  ))
  const actionExample = getActionExample(intentType, authConfig.authType)
  const vars = { intentName: name, intentType, actionExample }
  const inDir = path.join(__dirname, 'templates/generate/intent')
  const outDir = path.join(path.dirname(rootPathRc), 'intents')
  copy(inDir, outDir, vars, (err, createdFiles) => {
    if (err) throw err
    createdFiles.forEach(filePath =>
      emitter.emit('generateIntent:fileGenerated', filePath)
    )
  })
}

module.exports = {
  useWith: (program, emitter, config) => {
    program
      .command('generate')
      .alias('g')
      .description(
        `Generate intent or screen.
    $ bearer generate
`
      )
      // .option('-t, --type <intentType>', 'Intent type.')
      .option('--setup', 'generate setup file')
      .option('--config', 'generate config file')
      .action(generate(emitter, config))
  }
}
