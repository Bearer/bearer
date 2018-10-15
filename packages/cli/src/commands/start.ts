import { flags } from '@oclif/command'
import * as chokidar from 'chokidar'

import BaseLegacyCommand from '../BaseLegacyCommand'
import { Config } from '../types'
import Locator from '../utils/locator'
import setupConfig from '../utils/setupConfig'
import StartLocalDevelopmentServer from '../utils/StartLocalDevelopmentServer'

import BuildIntents from './build/intents'
import GenerateSetup from './generate/setup'
import GenerateSpec from './generate/spec'
import PrepareIntents from './prepare/intents'
import PrepareViews from './prepare/views'

const noOpen = 'no-open'
const noInstall = 'no-install'
const noWatcher = 'no-watcher'

export default class Start extends BaseLegacyCommand {
  static description = 'Start local development environment'


  get bearerConfig(): Config {
    const { flags } = this.parse(this.constructor as any)
    const path = flags.path || undefined
    return setupConfig(path)
  }

  get locator(): Locator {
    return new Locator(this.bearerConfig)
  }

  static flags = {
    help: flags.help({ char: 'h' }),
    path: flags.string({}),
    [noOpen]: flags.boolean({}),
    [noInstall]: flags.boolean({}),
    [noWatcher]: flags.boolean({ hidden: true })
  }

  static args = []

  async run() {
    const { flags } = this.parse(Start)
    const cmdArgs = []
    if (flags[noOpen]) {
      cmdArgs.push(`--${noOpen}`)
    }
    if (flags[noInstall]) {
      cmdArgs.push(`--${noInstall}`)
    }
    if (flags[noWatcher]) {
      cmdArgs.push(`--${noWatcher}`)
    }

    await GenerateSetup.run(['--silent'])
    await GenerateSpec.run(['--silent'])
    await PrepareViews.run(['--silent'])
    await PrepareIntents.run(['--silent'])
    await BuildIntents.run(['--silent'])

    chokidar
      .watch('.', {
        ignored: /(^|[\/\\])\../,
        cwd: this.locator.srcIntentsDir,
        ignoreInitial: true,
        persistent: true,
        followSymlinks: false
      })
      .on('add', async () => { await BuildIntents.run(['--silent']) })
      .on('change', async () => { await BuildIntents.run(['--silent']) })

    await StartLocalDevelopmentServer.run(this)
    //this.runLegacy(['start', ...cmdArgs])
  }
}
