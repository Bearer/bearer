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
    term.red('Failed pushing to developer portal.\n')
    term.yellow('error message: ').red(error)
    term('\n')
    outputError(error)
  })

  emitter.on('developerPortalUpdate:error', error => {
    term.red('There was an error while pushing to developer portal.')
    term('\n')
    outputError(error)
  })

  /* ********* Start output ********* */

  emitter.on('start:prepare:installingDependencies', () => {
    term.yellow('Installing views dependencies.')
    term('\n')
  })

  emitter.on('start:watchers', () => {
    term.yellow('Starting watchers')
    term('\n')
  })

  emitter.on('start:watchers:stdout', ({ name, data }) => {
    term.yellow(`[watcher:${name}] `)
    term.green(data)
  })

  emitter.on('start:watchers:stderr', ({ name, data }) => {
    term.yellow(`[watcher:${name}] `)
    term.green(data)
  })

  emitter.on('start:watchers:close', ({ name, code }) => {
    term.yellow(`[watcher:${name}] closed exit code: ${code}\n`)
  })

  emitter.on('start:prepare:failed', ({ error }) => {
    term.red('Prepare : An error occured')
    term('\n')
    term.white('    Error: ')
    term.red(error)
    term('\n')
  })

  emitter.on('start:failed', ({ error }) => {
    term.red('An error occured')
    term('\n')
    term.white('    Error: ')
    term.red(error)
    term('\n')
  })

  emitter.on('start:localServer:start', ({ port }) => {
    term.yellow('[local:functionServer] ')
    term.yellow('Serving: ')
    term(`http://127.0.0.1:${port}`)
    term('\n')
  })

  emitter.on('start:localServer:customPort', ({ port, host }) => {
    term('*************** Action required *****************\n')
    term.red('[local:functionServer] You have specified a custom port.\n')
    term.yellow('You must update the views/index.html to match this setting as follow:\n\n')
    term(`<script> bearer("CLIENT_ID", { integrationHost: "${host}" }) </script>\n\n`)
    term('*************************************************\n')
  })

  emitter.on('start:localServer:endpoints', ({ endpoints }) => {
    term.yellow('[local:functionServer] ')
    term.yellow('paths:')
    term('\n\t* ')
    term(endpoints.map(i => i.path).join('\n\t* '))
    term('\n')
  })

  emitter.on('start:localServer:generatingFunctions:start', () => {
    term.yellow('[local:functionServer] ')
    term.yellow('Reloading functions')
    term('\n')
  })

  emitter.on('start:localServer:generatingFunctions:stop', () => {
    term.yellow('[local:functionServer] ')
    term.yellow('Functions reloaded')
    term('\n')
  })

  emitter.on('start:localServer:generatingFunctions:failed', ({ error }) => {
    term.yellow('[local:functionServer] ')
    term.red('Functions building error\n')
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
