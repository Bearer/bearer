#!/usr/bin/env node
import Transpiler from '../src'

export default args => {
  const transpiler = new Transpiler({
    watchFiles: args.indexOf('--no-watcher') === -1
  })

  process.on('message', message => {
    if (message === 'refresh') {
      transpiler.refresh()
    }
  })

  transpiler.on('STOP', () => {
    process.exit(0)
  })

  transpiler.run()
}
