const term = require('terminal-kit').terminal

function outputError(error) {
  if (typeof error !== 'undefined' && error.message) {
    term.white('Error: ')
    term.red(error.message)
    term('\n')
  }
}

function inviteCommand(command) {
  const padding = 5
  const separator = command.length + 2 * padding
  term.white('Bearer: ').yellow('Please run the command below\n')
  term.white('='.repeat(separator))
  term('\n')
  term(' '.repeat(padding))
  term.white(command)
  term(' '.repeat(padding))
  term('\n')
  term.white('='.repeat(separator))
  term('\n')
}

module.exports = emitter => {
  emitter.on('buildArtifact:output:close', outputPath => {
    term.white('Bearer: ')
    term.yellow('Artifact stored in ')
    term(outputPath)
    term('\n')
  })

  emitter.on('buildArtifact:output:end', () => {
    console.log('Data has been drained')
  })

  emitter.on('buildArtifact:failed', ({ error }) => {
    term
      .white('Bearer: ')
      .red('Error while building scenario artifact\n')
      .white('Bearer: ')
      .yellow('Error message: ')
      .white(error)
    term('\n')
  })

  emitter.on('buildArtifact:archive:warning:ENOENT', err => {
    term.red(`Warning ${err}`)
  })

  emitter.on('buildArtifact:start', ({ scenarioUuid }) => {
    term.white('Bearer: ')
    term.yellow(`Building scenario package artifact... [${scenarioUuid}]`)
    term('\n')
  })

  emitter.on('buildArtifact:error', errors => {
    errors.forEach(line => console.log(line))
    term('\n')
  })
  emitter.on('buildArtifact:configured', ({ intents }) => {
    const intentNames = intents.map(intent => Object.keys(intent))
    term.white('Bearer: ')
    term.yellow(`Artefact configured : ${intentNames.join(' | ')}`)
    term('\n')
  })
  emitter.on('pushScenario:unauthorized', error => {
    outputError(error)
    inviteCommand('bearer login')
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
    term.red(`There was an error when trying to push the package: ${packagePath}\n`)
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

  emitter.on('assemblyScenario:success', () => {
    term.white('Bearer: ')
    term.yellow('Intents created.')
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

  emitter.on('generateTemplate:fileGenerated', path => {
    term.white('Bearer: ')
    term.yellow(`Bootstrapped a file: ${path}`)
    term('\n')
  })

  emitter.on('generateTemplate:deleteFiles', file => {
    term.white('Bearer: ')
    term.yellow(`Deleted a file: ${file}`)
    term('\n')
  })

  emitter.on('generateTemplate:skipped', configKey => {
    term.white('Bearer: ')
    term.red(`Skipped template generation: ${configKey} Not Found`)
    term('\n')
  })

  emitter.on('generateTemplate:error', error => {
    term.white('Bearer: ')
    term.red(`Error while generating template: ${error}`)
    term('\n')
  })

  emitter.on('generateView:fileGenerated', path => {
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
    term.red(`There was en error while trying to fetch authenticate user: ${JSON.stringify(err)}`)
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

  emitter.on('scenarioUuid:missing', devId => {
    term.white('Bearer: ')
    term.red('Missing scenarioUuid.\n')
    inviteCommand('bearer link')
    term('\n')
  })

  emitter.on('username:missing', () => {
    term.white('Bearer: ')
    term.red('Missing username.')
    term('\n')
    inviteCommand('bearer signup')
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
    inviteCommand('bearer new')
    term.yellow('to bootstrap a new scenario.')
    term('\n')
  })

  emitter.on('intents:installingDependencies', () => {
    term.white('Bearer: ')
    term.yellow('Installing intents dependencies.')
    term('\n')
  })

  emitter.on('views:generateSetupComponent', () => {
    term.white('Bearer: ')
    term.yellow('Generating setup component.')
    term('\n')
  })

  emitter.on('views:buildingDist', () => {
    term.white('Bearer: ')
    term.yellow('Building dist.')
    term('\n')
  })

  emitter.on('views:pushingDist', () => {
    term.white('Bearer: ')
    term.yellow('Pushing dist.')
    term('\n')
  })

  emitter.on('view:upload:start', () => {
    term.white('Bearer: ')
    term.yellow('Uploading views...')
    term('\n')
  })

  emitter.on('view:upload:success', () => {
    term.white('Bearer: ')
    term.yellow('Views uploaded successfully.')
    term('\n')
  })

  emitter.on('view:upload:error', e => {
    term.white('Bearer: ')
    term.red('ERROR: Views upload failed.')
    term('\n')
    term(e.toString())
    term('\n')
  })

  emitter.on('view:fileUpload:error', ({ e, key }) => {
    term.white('Bearer: ')
    term.red('ERROR: View file upload failed.')
    term('\n')
    term(`${key}: ${e.toString()}`)
    term('\n')
  })
  emitter.on('view:fileUpload:success', distPath => {
    term(distPath)
    term('\n')
  })

  emitter.on('view:fileUpload:failure', distPath => {
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

  emitter.on('user:notAuthenticated', () => {
    term.white('Bearer: ')
    term.red('There was an error while trying to retrieve your access token')
    term('\n')
    inviteCommand('bearer login')
    term('\n')
  })

  emitter.on('login:success', () => {
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
  emitter.on('login:failure', error => {
    term.white('Bearer: ')
    term.red('There was an error while trying to login to bearer')
    term('\n')
    outputError(error)
  })

  emitter.on('login:error', ({ body: error }) => {
    term.white('Bearer: ')
    term.red('There was an error while trying to login to bearer')
    term('\n')
    outputError(error)
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
    term.yellow('View invalidation success.')
    term('\n')
  })

  emitter.on('invalidateCloudFront:invalidationFailed', error => {
    term.white('Bearer: ')
    term.red("Couldn't invalidate views cache.")
    term('\n')
    outputError(error)
  })

  emitter.on('invalidateCloudFront:error', error => {
    term.white('Bearer: ')
    term.red('There was an error while trying to invalidate views cache.')
    term('\n')
    outputError(error)
  })

  emitter.on('developerPortalUpdate:failed', error => {
    term.white('Bearer: ')
    term.red('Failed pushing to developer portal.\n')
    term.white('Bearer: ')
    term.yellow('error message: ').red(error)
    term('\n')
    outputError(error)
  })

  emitter.on('developerPortalUpdate:error', error => {
    term.white('Bearer: ')
    term.red('There was an error while pushing to developer portal.')
    term('\n')
    outputError(error)
  })

  /* ********* Start output ********* */

  emitter.on('start:prepare:buildFolder', () => {
    term.white('Bearer: ')
    term.red('Generating build folder ')
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
    term.yellow('Installing views dependencies.')
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

  emitter.on('start:localServer:generatingIntents:start', () => {
    term.white('Bearer: ')
    term.yellow('[local:intentServer] ')
    term.yellow('Reloading intents')
    term('\n')
  })

  emitter.on('start:localServer:generatingIntents:stop', () => {
    term.white('Bearer: ')
    term.yellow('[local:intentServer] ')
    term.yellow('Intents reloaded')
    term('\n')
  })

  emitter.on('start:localServer:generatingIntents:failed', ({ error }) => {
    term.white('Bearer: ')
    term.yellow('[local:intentServer] ')
    term.red('Intents building error\n')
    console.log(
      error.toString({
        builtAt: false,
        entrypoints: false,
        assets: false,
        version: false,
        timings: false,
        hash: false,
        modules: false,
        chunks: false, // Makes the build much quieter
        colors: true // Shows colors in the console
      })
    )
    term('\n')
  })

  emitter.on('deployIntents:error', () => {
    term.white('Bearer: ')
    term.red('No intents found: Skipped')
    term('\n')
  })

  emitter.on('refreshToken:failure', () => {
    term.white('Bearer: ')
    term.red('Error while trying to authenticate\n')
    inviteCommand('bearer login')
    term('\n')
  })

  emitter.on('deployScenario:deployViews:error', error => {
    term.white('Bearer: ')
    term.red('An error occured')
    term('\n')
    outputError(error)
  })
  emitter.on('deployScenario:deployViews:error', e => {
    term.white('Bearer: ')
    term.red('An error occured')
    term('\n')
    term.white('Error: ')
    term.red(e.toString())
    term('\n')
  })
}
