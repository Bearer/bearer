const term = require('terminal-kit').terminal

module.exports = emitter => {
  emitter.on('buildArtifact:output:close', output_path => {
    term.white('Bearer: ')
    term.yellow('Artifact stored in ')
    term(output_path)
    term('\n')
  })

  emitter.on('buildArtifact:output:end', () => {
    console.log('Data has been drained')
  })

  emitter.on('buildArtifact:archive:warning:ENOENT', err => {
    term.red(`Warning ${err}`)
  })

  emitter.on('buildArtifact:start', ({ scenarioUuid }) => {
    term.white('Bearer: ')
    term.yellow(`Building scenario package artifact... [${scenarioUuid}]`)
    term('\n')
  })

  emitter.on('buildArtifact:configured', ({ intents }) => {
    const intentNames = intents.map(intent => Object.keys(intent))
    term.white('Bearer: ')
    term.yellow(`Artefact configured : ${intentNames.join(' | ')}`)
    term('\n')
  })
  emitter.on('pushScenario:unauthorized', ({ message }) => {
    term.white('Bearer: ')
    term.red(`ERROR: ${message}`)
    term('\n')
    term.white('Bearer: ')
    term.red(`Please try to `)
    term('bearer login ')
    term.red('again.')
    term('\n')
  })

  emitter.on('pushScenario:httpError', res => {
    term.white('Bearer: ')
    term.red(`ERROR: ${JSON.stringify(res, null, 2)}`)
    term('\n')
  })

  emitter.on('pushScenario:error', error => {
    term.white('Bearer: ')
    term.red(`There was an error when trying to push scenario: `)
    term('\n')
    term(error.toString())
    term('\n')
  })

  emitter.on('pushScenario:start', Key => {
    term.white('Bearer: ')
    term.yellow(`Pushing scenario ${Key}...`)
    term('\n')
  })

  emitter.on('pushScenario:uploadPackage:error', (err, packagePath) => {
    term.red(
      `There was an error when trying to push the package: ${packagePath}\n`
    )
    term.white(`Error: ${err}\n`)
  })

  emitter.on('pushScenario:uploadPackage:success', () => {
    term.white('Bearer: ')
    term.yellow('Scenario has been uploaded.')
    term('\n')
  })

  emitter.on('assemblyScenario:start', () => {
    term.white('Bearer: ')
    term.yellow('Building intents...')
    term('\n')
  })

  emitter.on('assemblyScenario:success', body => {
    term.white('Bearer: ')
    term.yellow('Intents created.')
    term('\n')

    term.white('PackageManager: ')
    term.yellow(JSON.stringify(body))
    term('\n')
  })

  emitter.on('assemblyScenario:failed', response => {
    term.white('Bearer: ')
    term.yellow('Something went wrong...')
    term(JSON.stringify(response, null, 2))
    term('\n')
  })

  emitter.on('assemblyScenario:error', err => {
    term.white('Bearer: ')
    term.red(`There was an error while trying to start the assembly: ${err}`)
    term('\n')
  })

  emitter.on('generateIntent:fileGenerated', path => {
    term.white('Bearer: ')
    term.yellow(`Bootstrapped a file: ${path}`)
    term('\n')
  })
  emitter.on('signUp:userCreated', email => {
    term.white('Bearer: ')
    term.yellow('User created: ')
    term.white(email)
    term('\n')
  })

  emitter.on('credentialsUpdated', configPath => {
    term.white('Bearer: ')
    term.yellow('Credentials and configuration stored in: ')
    term.white(configPath)
    term('\n')
  })
  emitter.on('signin:authenticateUser:getUserAttirbutes:failed', err => {
    term.white('Bearer: ')
    term.red(`There was en error while trying to fetch user attributes: ${err}`)
    term('\n')
  })

  emitter.on('signin:authenticateUser:failed', err => {
    term.white('Bearer: ')
    term.red(
      `There was en error while trying to fetch authenticate user: ${JSON.stringify(
        err
      )}`
    )
    term('\n')
  })

  emitter.on('signUp:error', err => {
    term.white('Bearer: ')
    term.red(`There was en error while trying to fetch sign up an user: ${err}`)
    term('\n')
  })

  emitter.on('developerIdFound', devId => {
    term.white('Bearer: ')
    term.red(`Your developerId: ${devId}`)
    term('\n')
  })

  emitter.on('scenarioTitle:missing', devId => {
    term.white('Bearer: ')
    term.red('Missing scenarioTitle')
    term('\n')
  })

  emitter.on('username:missing', () => {
    term.white('Bearer: ')
    term.red('Missing username.')
    term('\n')
    term.white('Bearer: ')
    term.yellow('Run ')
    term('bearer signup --email <email>')
    term.yellow(' first.')
    term('\n')
  })
  emitter.on('scenarioTitle:creationFailed', e => {
    term.white('Bearer: ')
    term.red("Couldn't store the scenarioTitle")
    term('\n')
    term(e)
    term('\n')
  })
  emitter.on('rootPath:doesntExist', () => {
    term.white('Bearer: ')
    term.red('Looks like you are not in scenario project directory.')
    term('\n')
    term.white('Bearer: ')
    term.red('Run ')
    term('bearer new <scenarioTitle>')
    term.red(' to bootstrap a new scenario.')
    term('\n')
  })

  emitter.on('intents:installingDependencies', () => {
    term.white('Bearer: ')
    term.yellow('Installing intents dependencies.')
    term('\n')
  })

  emitter.on('screens:generateSetupComponent', () => {
    term.white('Bearer: ')
    term.yellow('Generating setup component.')
    term('\n')
  })

  emitter.on('screens:buildingDist', () => {
    term.white('Bearer: ')
    term.yellow('Building dist.')
    term('\n')
  })

  emitter.on('screens:pushingDist', () => {
    term.white('Bearer: ')
    term.yellow('Pushing dist.')
    term('\n')
  })

  emitter.on('screen:upload:start', () => {
    term.white('Bearer: ')
    term.yellow('Uploading screens...')
    term('\n')
  })

  emitter.on('screen:upload:success', () => {
    term.white('Bearer: ')
    term.yellow('Screens uploaded successfully.')
    term('\n')
  })

  emitter.on('screen:upload:error', e => {
    term.white('Bearer: ')
    term.red('ERROR: Screens upload failed.')
    term('\n')
    term(e.toString())
    term('\n')
  })

  emitter.on('screen:fileUpload:error', e => {
    term.white('Bearer: ')
    term.red('ERROR: Screen file upload failed.')
    term('\n')
    term(e.toString())
    term('\n')
  })
  emitter.on('screen:fileUpload:success', distPath => {
    term(distPath)
    term('\n')
  })

  emitter.on('screen:fileUpload:failure', distPath => {
    term.white('Bearer: ')
    term.red("Couldn't upload a file")
    term('\n')
    term(distPath)
    term('\n')
  })

  emitter.on('storeCredentials:success', referenceId => {
    term.white('Bearer: ')
    term.yellow('referenceId: ')
    term(referenceId)
    term('\n')
  })

  emitter.on('storeCredentials:failure', e => {
    term.white('Bearer: ')
    term.red('There was an error while trying to save credentials')
    term('\n')
    term(e.toString())
    term('\n')
  })

  emitter.on('signup:success', body => {
    term.white('Bearer: ')
    term.yellow('successfully signed up to bearer.')
    term('\n')
    term.white('Bearer: ')
    term(JSON.stringify(body))
    term.yellow(' saved to ')
    term('~/.bearerrc')
    term('\n')
  })

  emitter.on('signup:failure', body => {
    term.white('Bearer: ')
    term.red('There was an error while trying to signup to bearer')
    term('\n')
    term.white('IntegrationService: ')
    term.red(body.message)
    term('\n')
  })

  emitter.on('signup:error', e => {
    term.white('Bearer: ')
    term.red('There was an error while trying to signup to bearer')
    term('\n')
    term.white('Error: ')
    term.red(e.toString())
    term('\n')
  })

  emitter.on('login:success', body => {
    term.white('Bearer: ')
    term.yellow('successfully logged in to bearer.')
    term('\n')
    term.white('Bearer: ')
    term.yellow('AccessToken saved to ')
    term('~/.bearerrc')
    term('\n')
  })

  emitter.on('login:userFound', Username => {
    term.white('Bearer: ')
    term.yellow('Hello ')
    term(Username)
    term.yellow('!')
    term('\n')
  })
  emitter.on('login:failure', ({ message }) => {
    term.white('Bearer: ')
    term.red('There was an error while trying to login to bearer')
    term('\n')
    term.white('IntegrationService: ')
    term.red(message)
    term('\n')
  })

  emitter.on('login:error', ({ body: { message } }) => {
    term.white('Bearer: ')
    term.red('There was an error while trying to login to bearer')
    term('\n')
    term.white('Error: ')
    term.red(message)
    term('\n')
  })

  emitter.on('deploy:started', () => {
    term.white('Bearer: ')
    term.yellow('Starting scenario deployment')
    term('\n')
  })

  emitter.on('deploy:finished', ({ setupUrl }) => {
    term.white('Bearer: ')
    term('\n')
    term.yellow(`Scenario setup: `)
    term.white(setupUrl)
    term('\n')
  })

  emitter.on('invalidateCloudFront:success', () => {
    term.white('Bearer: ')
    term.yellow('Screen invalidation success.')
    term('\n')
  })

  emitter.on('invalidateCloudFront:invalidationFailed', ({ message }) => {
    term.white('Bearer: ')
    term.red("Couldn't invalidate screens cache.")
    term('\n')
    term.white('Error: ')
    term.red(message)
    term('\n')
  })

  emitter.on('invalidateCloudFront:error', ({ message }) => {
    term.white('Bearer: ')
    term.red('There was an error while trying to invalidate screens cache.')
    term('\n')
    term.white('Error: ')
    term.red(message)
    term('\n')
  })

  /* ********* Start output ********* */

  emitter.on('start:prepare:buildFolder', () => {
    term.white('Bearer: ')
    term.yellow('Generating .build folder ')
    term('\n')
  })

  emitter.on('start:prepare:stencilConfig', () => {
    term.white('Bearer: ')
    term.yellow('Generating stencil configuration')
    term('\n')
  })

  emitter.on('start:prepare:copyFile', file => {
    term.white('Bearer: ')
    term.yellow(`Copied: ${file}`)
    term('\n')
  })

  emitter.on('start:symlinkNodeModules', () => {
    term.white('Bearer: ')
    term.yellow('Symlinking node_modules')
    term('\n')
  })

  emitter.on('start:symlinkPackage', () => {
    term.white('Bearer: ')
    term.yellow('Symlinking package.json')
    term('\n')
  })

  emitter.on('start:prepare:installingDependencies', () => {
    term.white('Bearer: ')
    term.yellow('Installing screens dependencies.')
    term('\n')
  })

  emitter.on('start:watchers', () => {
    term.white('Bearer: ')
    term.yellow('Starting watchers')
    term('\n')
  })

  emitter.on('start:watchers:stdout', ({ name, data }) => {
    term.white('Bearer: ')
    term.yellow(`[watcher:${name}] `)
    term.green(data)
  })

  emitter.on('start:watchers:stderr', ({ name, data }) => {
    term.white('Bearer: ')
    term.yellow(`[watcher:${name}] `)
    term.green(data)
  })

  emitter.on('start:watchers:close', ({ name, code }) => {
    term.white('Bearer: ')
    term.yellow(`[watcher:${name}] closed exit code: ${code}\n`)
  })

  // emitter.emit('start:watchers:stencil:stdout', )

  emitter.on('start:prepare:failed', ({ error }) => {
    term.white('Bearer: ')
    term.red('Prepare : An error occured')
    term('\n')
    term.white('    Error: ')
    term.red(error)
    term('\n')
  })

  emitter.on('start:failed', ({ error }) => {
    term.white('Bearer: ')
    term.red('An error occured')
    term('\n')
    term.white('    Error: ')
    term.red(error)
    term('\n')
  })

  emitter.on('start:localServer:start', ({ port }) => {
    term.white('Bearer: ')
    term.yellow('[local:intentServer] ')
    term.yellow('Serving: ')
    term(`http://127.0.0.1:${port}`)
    term('\n')
  })

  emitter.on('start:localServer:endpoints', ({ endpoints }) => {
    term.white('Bearer: ')
    term.yellow('[local:intentServer] ')
    term.yellow('paths: ')
    term(endpoints.map(i => i.path).join(', '))
    term('\n')
  })

  emitter.on('deployScenario:deployScreens:error', ({ message }) => {
    term.white('Bearer: ')
    term.red('An error occured')
    term('\n')
    term.white('Error: ')
    term.red(message)
    term('\n')
  })
}
