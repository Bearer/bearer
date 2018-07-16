#!/usr/bin/env node
import Transpiler, * as Trasnpiler from '../src/index'

const args = process.argv.slice(2)

const transpiler = new Transpiler(
  undefined,
  args.indexOf('--no-watcher') === -1
)

transpiler.on('STOP', () => {
  process.exit(0)
})

transpiler.run()
