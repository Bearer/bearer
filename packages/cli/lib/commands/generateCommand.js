const copy = require('copy-template-dir')
const del = require('del')
const path = require('path')
const inquirer = require('inquirer')
const Case = require('case')
const fs = require('fs')
const util = require('util')
const intents = require('@bearer/intents')

const INTENT = 'intent'
const SCREEN = 'screen'
const COLLECTION = 'GetCollection'
const SINGLE_RESOURCE = 'GetObject'

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

  const { template, name } = await inquirer.prompt([
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
    },
    {
      message: 'Give it a name',
      type: 'input',
      name: 'name'
    }
  ])

  const params = { emitter, rootPathRc, name }

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

async function generateTemplates({ emitter, templateType, rootPathRc }) {
  const authConfig = require(path.dirname(rootPathRc) +
    '/intents/auth.config.json')
  const scenariorc = fs.readFileSync(path.dirname(rootPathRc) + '/.scenariorc')
  const scenarioId = scenariorc
    .toString()
    .split('=')[1]
    .replace(/[\n\r]+/g, '')
    .trim()

  const vars = {
    scenarioTitle: Case.camel(scenarioId),
    componentTagName: Case.kebab(scenarioId),
    fields: authConfig[templateType]
      ? JSON.stringify(authConfig[templateType])
      : '[]'
  }
  const inDir = path.join(__dirname, `templates/generate/${templateType}`)
  const outDir = path.join(path.dirname(rootPathRc), '/screens/src/')

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

function generateScreen({ emitter, rootPathRc, name }) {
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

const choices = Object.keys(intents).map(intent => ({
  name: intent,
  value: intent
}))

async function generateIntent({ emitter, rootPathRc, name }) {
  const { intentType } = await inquirer.prompt([
    {
      message: 'What type of intent do you wan to generate',
      type: 'list',
      name: 'intentType',
      choices
    }
  ])
  const vars = { intentName: name, intentType }
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
