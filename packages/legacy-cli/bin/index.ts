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
const generateCmd = require('../src/lib/commands/generateCommand')
const initCmd = require('../src/lib/commands/initCommand')
const loginCmd = require('../src/lib/commands/loginCommand')
const startCmd = require('../src/lib/commands/startCommand')
const linkCmd = require('../src/lib/commands/linkCommand')
const invokeCmd = require('../src/lib/commands/invokeCommand')

const cliOutput = require('../src/lib/cliOutput.js')

const cli = new CLI(program, emitter, config)
cliOutput(emitter, config)

program.version(version, '-v, --version')

cli.use(initCmd)
cli.use(generateCmd)
cli.use(deployCmd)
cli.use(loginCmd)
cli.use(startCmd)
cli.use(linkCmd)
cli.use(invokeCmd)

cli.parse(process.argv)
