#!/usr/bin/env node

require('../scripts/check-version')
const { version } = require('../package.json')
const { CLI } = require('../lib/cli')
const setupConfig = require('../lib/setupConfig')
const Emitter = require('../lib/emitter')

const emitter = new Emitter()
const inquirer = require('inquirer')

const config = setupConfig()
const program = require('commander')
const deployCmd = require('../lib/commands/deployCommand')
const generateCmd = require('../lib/commands/generateCommand')
const initCmd = require('../lib/commands/initCommand')
const loginCmd = require('../lib/commands/loginCommand')

const startCmd = require('../lib/commands/startCommand')

const cliOutput = require('../lib/cliOutput.js')

const cli = new CLI(program, emitter, config)
cliOutput(emitter, config)

program.version(version, '-v, --version')

cli.use(initCmd)
cli.use(generateCmd)
cli.use(deployCmd)
cli.use(loginCmd)
cli.use(startCmd)

cli.parse(process.argv)
