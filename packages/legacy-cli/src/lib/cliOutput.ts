import debug from './logger'
const logger = debug.extend('out')

const term = require('terminal-kit').terminal

function outputError(error) {
  if (typeof error !== 'undefined' && error.message) {
    term.white('Error: ')
    term.red(error.message)
    term('\n')
  }
}

export default emitter => {
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
    term.yellow('paths:')
    term('\n\t* ')
    term(endpoints.map(i => i.path).join('\n\t* '))
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
    logger(
      '%s',
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
