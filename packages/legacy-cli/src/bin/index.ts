#!/usr/bin/env node
const program = require('commander')

// Done at OCLIF level
// require('../scripts/check-version')
const { version } = require('../../package.json')
const { CLI } = require('../lib/cli')
// tslint:disable-next-line
const Emitter = require('../lib/emitter')

import setupConfig from '../lib/setupConfig'
import { Config } from '../lib/types'

const emitter = new Emitter()
const config: Config = setupConfig()
const startCmd = require('../src/lib/commands/startCommand')
const invokeCmd = require('../src/lib/commands/invokeCommand')

const cliOutput = require('../src/lib/cliOutput.js')

const cli = new CLI(program, emitter, config)
cliOutput(emitter, config)

program.version(version, '-v, --version')

cli.use(startCmd)
cli.use(invokeCmd)

export default args => {
  cli.parse(args)
}
