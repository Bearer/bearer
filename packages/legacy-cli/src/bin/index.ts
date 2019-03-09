#!/usr/bin/env node
import * as program from 'commander'

// Done at OCLIF level
// require('../scripts/check-version')
import cliOutput from '../lib/cliOutput'
import { CLI } from '../lib/cli'
import Emitter from '../lib/emitter'
import setupConfig from '../lib/setupConfig'
import { Config } from '../lib/types'
import * as startCmd from '../lib/commands/startCommand'

export default args => {
  const emitter = new Emitter()
  const config: Config = setupConfig()

  const cli = new CLI(program, emitter, config)
  cliOutput(emitter)
  cli.use(startCmd)

  cli.parse(args)
}
