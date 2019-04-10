#!/usr/bin/env node
import commandLineArgs from 'command-line-args'

import Transpiler from '../src'

export default args => {
  const optionsDefinitions = [
    // CORE-227: see later in the file
    { name: 'build', type: Boolean, defaultValue: false },
    { name: 'no-watcher', type: Boolean, defaultValue: false },
    { name: 'no-process', type: Boolean, defaultValue: false },
    { name: 'prefix-tag', type: String },
    { name: 'suffix-tag', type: String },
    { name: 'buid', type: String },
    { name: 'fail-fast', type: Boolean, defaultValue: false }
  ]

  const options = commandLineArgs(optionsDefinitions, {
    camelCase: true,
    argv: args,
    partial: true
  })
  const transpiler = new Transpiler({
    watchFiles: !options.noWatcher,
    tagNamePrefix: options.prefixTag,
    tagNameSuffix: options.suffixTag,
    failFast: options.failFast,
    buid: options.buid
  })
  if (!options.noProcess) {
    process.on('message', message => {
      if (message === 'refresh') {
        transpiler.refresh()
      }
    })

    transpiler.on('STOP', () => {
      process.exit(0)
    })
  }

  transpiler.on('ERROR', (data: { error: any }) => {
    console.log('Error while transpiling: %j', data.error.message)
    process.exit(1)
  })

  transpiler.run()
  // hack: tell the tranpiler to refresh and ensure metadata is up to date
  // CORE-227
  if (options.build) {
    transpiler.refresh()
  }
}
