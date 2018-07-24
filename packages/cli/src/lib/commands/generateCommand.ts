import * as copy from 'copy-template-dir'
import * as del from 'del'
import * as path from 'path'
import * as inquirer from 'inquirer'
import * as Case from 'case'
import * as intents from '@bearer/intents'
import * as templates from '@bearer/templates'
import * as rc from 'rc'
import Locator from '../locationProvider'

const INTENT = 'intent'
const SCREEN = 'screen'
enum TemplateTypes {
  config = 'config',
  setup = 'setup'
}
async function generateTemplates({
  emitter,
  templateType,
  locator
}: {
  emitter: any
  templateType: TemplateTypes
  locator: Locator
}) {
  const authConfig = require(locator.scenarioRootFile('auth.config.json'))

  const scenarioConfig = rc('scenario')
  const { scenarioTitle } = scenarioConfig

  const configKey = `${templateType}Screens`

  const inDir = path.join(__dirname, `templates/generate/${templateType}`)
  const outDir = locator.buildScreenDir

  await del(`${outDir}*${templateType}*.tsx`).then(paths => {
    console.log('Deleted files and folders:\n', paths.join('\n'))
  })

  if (authConfig[configKey] && authConfig[configKey].length) {
    const vars = {
      scenarioTitle: Case.camel(scenarioTitle),
      componentTagName: Case.kebab(scenarioTitle),
      fields: JSON.stringify(authConfig[configKey])
    }

    copy(inDir, outDir, vars, (err, createdFiles) => {
      if (err) throw err
      createdFiles.forEach(filePath => emitter.emit('generateIntent:fileGenerated', filePath))
    })
  }
}

const generate = (emitter, {}, locator: Locator) => async env => {
  const { scenarioRoot } = locator
  if (!scenarioRoot) {
    emitter.emit('rootPath:doesntExist')
    process.exit(1)
  }

  if (env.config) {
    return generateTemplates({
      emitter,
      templateType: TemplateTypes.config,
      locator
    })
  }

  if (env.setup) {
    return generateTemplates({
      emitter,
      templateType: TemplateTypes.setup,
      locator
    })
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

  const params = { emitter, locator }

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

async function generateScreen({ emitter, locator }: { locator: Locator; emitter: any }) {
  const name = await askForName()
  const componentName = Case.pascal(name)
  const vars = {
    screenName: componentName,
    componentTagName: Case.kebab(componentName)
  }
  const inDir = path.join(__dirname, 'templates/generate/screen')
  const outDir = path.join(locator.srcScreenDir, 'components')

  copy(inDir, outDir, vars, (err, createdFiles) => {
    if (err) throw err
    createdFiles.forEach(filePath => emitter.emit('generateScreen:fileGenerated', filePath))
  })
}

const filteredChoices = (intents: Record<string, any>, propertyFlag) =>
  Object.keys(intents)
    .filter(intent => intents[intent][propertyFlag])
    .map(intent => ({
      name: intents[intent].display,
      value: intent
    }))
    .sort((a, b) => (a.name > b.name ? 1 : -1))

const choices = [
  ...filteredChoices(intents, 'isGlobalIntent'),
  new inquirer.Separator(),
  ...filteredChoices(intents, 'isStateIntent')
]

function getActionExample(intentType, authType) {
  return templates[authType][intentType]
}

async function generateIntent({ emitter, locator }: { emitter: any; locator: Locator }) {
  const { intentType } = await inquirer.prompt([
    {
      message: 'What type of intent do you want to generate',
      type: 'list',
      name: 'intentType',
      choices
    }
  ])
  const name = await askForName()
  const authConfig = require(locator.scenarioRootFile('auth.config.json'))
  const actionExample = getActionExample(intentType, authConfig.authType)
  const vars = { intentName: name, authType: authConfig.authType, intentType, actionExample }
  const inDir = path.join(__dirname, 'templates/generate/intent')
  const outDir = locator.srcIntentDir

  copy(inDir, outDir, vars, (err, createdFiles) => {
    if (err) throw err
    createdFiles.forEach(filePath => emitter.emit('generateIntent:fileGenerated', filePath))
  })
}

export function useWith(program, emitter, config, locator): void {
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
    .action(generate(emitter, config, locator))
}
