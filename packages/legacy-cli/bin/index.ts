#!/usr/bin/env node
const program = require('commander')

// Done at OCLIF level
// require('../scripts/check-version')
const { version } = require('../../package.json')
const { CLI } = require('../src/lib/cli')
const Emitter = require('../src/lib/emitter')

import setupConfig from '../src/lib/setupConfig'
import { Config } from '../src/lib/types'

const emitter = new Emitter()
const config: Config = setupConfig()
const deployCmd = require('../src/lib/commands/deployCommand')
const startCmd = require('../src/lib/commands/startCommand')
const invokeCmd = require('../src/lib/commands/invokeCommand')

const cliOutput = require('../src/lib/cliOutput.js')

const cli = new CLI(program, emitter, config)
cliOutput(emitter, config)

program.version(version, '-v, --version')

cli.use(deployCmd)
cli.use(startCmd)
cli.use(invokeCmd)

export default args => {
  cli.parse(args)
}
