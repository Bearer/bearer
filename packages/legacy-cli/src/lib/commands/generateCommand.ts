import * as intents from '@bearer/intents'
import * as templates from '@bearer/templates'
import * as Case from 'case'
import * as copy from 'copy-template-dir'
import * as inquirer from 'inquirer'
import * as path from 'path'

import Locator from '../locationProvider'

import { generateSetup } from './generate'

const INTENT = 'intent'
const COMPONENT = 'component'

const enum Components {
  BLANK = 'blank',
  COLLECTION = 'collection',
  ROOT = 'root'
}

const generate = (emitter, {}, locator: Locator) => async env => {
  const { scenarioRoot } = locator
  if (!scenarioRoot) {
    emitter.emit('rootPath:doesntExist')
    process.exit(1)
  }

  if (env.setup) {
    return generateSetup({ emitter, locator })
  }

  if (env.blankComponent && typeof env.blankComponent === 'string') {
    return generateComponent({ emitter, locator, name: env.blankComponent, type: Components.BLANK })
  }

  if (env.blankComponent) {
    return generateComponent({ emitter, locator, type: Components.BLANK })
  }

  if (env.collectionComponent && typeof env.collectionComponent === 'string') {
    return generateComponent({ emitter, locator, name: env.collectionComponent, type: Components.COLLECTION })
  }

  if (env.collectionComponent) {
    return generateComponent({ emitter, locator, type: Components.COLLECTION })
  }

  if (env.rootGroup && typeof env.rootGroup === 'string') {
    return generateComponent({ emitter, locator, name: env.rootGroup, type: Components.ROOT })
  }

  if (env.rootGroup) {
    return generateComponent({ emitter, locator, type: Components.ROOT })
  }

  const { template } = await inquirer.prompt([
    {
      message: 'What do you want to generate',
      type: 'list',
      name: 'template',
      choices: [{ name: 'Intent', value: INTENT }, { name: 'Component', value: COMPONENT }]
    }
  ])

  const params = { emitter, locator }

  switch (template) {
    case INTENT:
      await generateIntent(params)
      break
    case COMPONENT:
      await generateComponent(params)
      break
    default:
  }
}

type TgenerateComponent = { locator: Locator; emitter: any; name?: string; type?: string }
async function generateComponent({ emitter, locator, name, type }: TgenerateComponent) {
  // Ask for type if not present
  if (!type) {
    type = await askComponentType()
  }

  // Ask for name if not present
  if (!name) {
    name = await askForName()
  }

  const inDir = path.join(__dirname, 'templates/generate', `${type}Component`)
  const outDir = type === Components.ROOT ? locator.srcViewsDir : locator.srcViewsDirResource('components')

  copy(inDir, outDir, getComponentVars(name, require(locator.authConfigPath)), (err, createdFiles) => {
    if (err) throw err
    createdFiles.forEach(filePath => emitter.emit('generateView:fileGenerated', filePath))
  })
}

async function generateIntent({ emitter, locator }: { emitter: any; locator: Locator }) {
  const intentType = await askIntentType()
  const name = await askForName()
  const authConfig = require(locator.authConfigPath)
  const inDir = path.join(__dirname, 'templates/generate/intent')
  const outDir = locator.srcIntentsDir

  copy(inDir, outDir, getIntentVars(name, intentType, authConfig), (err, createdFiles) => {
    if (err) throw err
    createdFiles.forEach(filePath => emitter.emit('generateTemplate:fileGenerated', filePath))
  })
}

/**
 * Helpers
 */

export function getIntentChoices(): Array<{ name: string; value: any }> {
  const filteredChoices = (intents: Record<string, any>, propertyFlag) =>
    Object.keys(intents)
      .filter(intent => intents[intent][propertyFlag])
      .map(intent => ({
        name: intents[intent].display,
        value: intent
      }))
      .sort((a, b) => (a.name > b.name ? 1 : -1))

  return [
    ...filteredChoices(intents, 'isGlobalIntent'),
    new inquirer.Separator(),
    ...filteredChoices(intents, 'isStateIntent')
  ]
}

function getActionExample(intentType, authType): string {
  return templates[authType][intentType]
}

async function askIntentType(): Promise<string> {
  const { intentType } = await inquirer.prompt([
    {
      message: 'What type of intent do you want to generate',
      type: 'list',
      name: 'intentType',
      choices: getIntentChoices()
    }
  ])
  return intentType
}

async function askComponentType(): Promise<Components> {
  const typePrompt = await inquirer.prompt([
    {
      message: 'What type of component do you want to generate',
      type: 'list',
      name: 'type',
      choices: [
        { name: 'Blank', value: Components.BLANK },
        { name: 'Collection', value: Components.COLLECTION },
        new inquirer.Separator(),
        { name: 'Root Group', value: Components.ROOT }
      ]
    }
  ])
  return typePrompt.type
}

export function getComponentVars(name: string, authConfig: { authType: string }) {
  const componentName = Case.pascal(name)
  return {
    fileName: name,
    componentName,
    componentClassName: componentName, // it gives more meaning within templates
    componentTagName: Case.kebab(componentName),
    groupName: Case.kebab(componentName),
    withAuthScreen: authConfig.authType !== 'noAuth' ? '<bearer-navigator-auth-screen />' : null
  }
}

export function getIntentVars(name: string, intentType: string, authConfig: { authType: string }) {
  const actionExample = getActionExample(intentType, authConfig.authType)
  return {
    fileName: name,
    intentName: name,
    intentClassName: Case.pascal(name),
    authType: authConfig.authType,
    intentType,
    actionExample,
    callbackType: `T${intentType}Callback`
  }
}

async function askForName() {
  const { name } = await inquirer.prompt([{ message: 'Give it a name', type: 'input', name: 'name' }])
  return name.trim()
}

/**
 * Command
 */
export function useWith(program, emitter, config, locator): void {
  program
    .command('generate')
    .alias('g')
    .description(
      `Generate intent or component.
    $ bearer generate
  `
    )
    // .option('-t, --type <intentType>', 'Intent type.')
    .option('--blank-component [name]', 'generate blank component')
    .option('--collection-component [name]', 'generate collection component')
    .option('--root-group [name]', 'generate root components group')
    .option('--setup', 'generate setup file')
    .action(generate(emitter, config, locator))
}
