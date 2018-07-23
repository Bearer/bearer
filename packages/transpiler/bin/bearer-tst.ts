#!/usr/bin/env node
import Transpiler from '../src/index'

export default args => {
  const transpiler = new Transpiler({
    watchFiles: args.indexOf('--no-watcher') === -1
  })

  transpiler.on('STOP', () => {
    process.exit(0)
  })

  transpiler.run()
}
