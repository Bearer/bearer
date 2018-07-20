#!/usr/bin/env node

require('../scripts/check-version')
const { version } = require('../../package.json')
const { CLI } = require('../src/lib/cli')
const setupConfig = require('../src/lib/setupConfig')
const Emitter = require('../src/lib/emitter')

const emitter = new Emitter()
const config = setupConfig()
const program = require('commander')
const deployCmd = require('../src/lib/commands/deployCommand')
const generateCmd = require('../src/lib/commands/generateCommand')
const initCmd = require('../src/lib/commands/initCommand')
const loginCmd = require('../src/lib/commands/loginCommand')

const startCmd = require('../src/lib/commands/startCommand')

const cliOutput = require('../src/lib/cliOutput.js')

const cli = new CLI(program, emitter, config)
cliOutput(emitter, config)

program.version(version, '-v, --version')

cli.use(initCmd)
cli.use(generateCmd)
cli.use(deployCmd)
cli.use(loginCmd)
cli.use(startCmd)

cli.parse(process.argv)
