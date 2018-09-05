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
  emitter.on('credentialsUpdated', configPath => {
    term.white('Bearer: ')
    term.yellow('Credentials and configuration stored in: ')
    term.white(configPath)
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

  emitter.on('user:notAuthenticated', () => {
    term.white('Bearer: ')
    term.red('There was an error while trying to retrieve your access token')
    term('\n')
    inviteCommand('bearer login')
    term('\n')
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
}
