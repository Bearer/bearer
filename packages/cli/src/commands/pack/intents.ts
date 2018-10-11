import prepareConfig from '@bearer/bearer-cli/src/lib/buildArtifact/prepareConfig'
import * as archiver from 'archiver'
import * as fs from 'fs-extra'
import * as path from 'path'

import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'

const CONFIG_FILE = 'bearer.config.json'
const HANDLER_NAME = 'index.js'

export default class PackIntents extends BaseCommand {
  static description = 'Pack scenario intents'
  static hidden = true
  static flags = {
    ...BaseCommand.flags
  }

  static args = [{ name: 'ARCHIVE_PATH', required: true }]

  @RequireScenarioFolder()
  async run() {
    // we are assuming prepare and build have been run previously
    const { args } = this.parse(PackIntents)
    const { ARCHIVE_PATH } = args
    const target = path.resolve(ARCHIVE_PATH)
    this.debug(`Zipping to: ${target}`)

    const output = fs.createWriteStream(target)
    const archive = archiver('zip', {
      zlib: { level: 9 }
    })
    const archiveEntries: Array<archiver.EntryData> = []
    archive.pipe(output)

    const archived = new Promise<Array<archiver.EntryData>>((resolve, reject) => {
      output.on('close', () => {
        resolve(archiveEntries)
      })
      archive.on('error', (err: any) => {
        reject(err)
      })

      archive.on('entry', (entry: archiver.EntryData) => {
        if (!entry.stats || !entry.stats.isDirectory()) {
          archiveEntries.push(entry)
        }
      })

      archive.on('warning', (err: any) => {
        this.debug(err)
      })
    })

    try {
      // add intents
      archive.directory(this.locator.buildIntentsDir, false)
      // add CONFIG
      const config = await this.appendConfig()
      archive.append(JSON.stringify(config, null, 2), { name: CONFIG_FILE })

      // add handler
      this.debug(`Generated config: ${JSON.stringify(config, null, 2)}`)
      archive.append(this.handlerContent(config), { name: HANDLER_NAME })
      // ZIP

      archive.finalize()

      const entries = await archived

      this.debug(`Zip content: ${entries.length} files\n  * ${entries.map(e => e.name).join('\n  * ')}`)
      this.success('Successfully generated lambda package')
      // log files added to zip
    } catch (e) {
      this.error(e)
    }
  }

  async appendConfig(): Promise<TConfig> {
    // generate config
    const config: TConfig = { intents: await this.retrieveIntents() }

    const content = await fs.readFile(this.locator.authConfigPath, { encoding: 'utf8' })
    config.auth = JSON.parse(content)
    return config
  }

  handlerContent({ intents }: TConfig): string {
    return intents
      .map(Object.keys)
      .map(
        intent => `
const ${intent} = require("./dist/${intent}").default;
module.exports[${intent}.intentName] = ${intent}.intentType.intent(${intent}.action);
`
      )
      .join('\n')
  }

  // TODO: rewrite this using TS AST
  async retrieveIntents(): Promise<Array<TIntentConfig>> {
    try {
      const config = await prepareConfig(this.locator.authConfigPath, "", this.bearerConfig.scenarioUuid, "", this.locator.srcIntentsDir)
      return config.intents
    } catch (e) {
      throw e
    }
  }
}
type TIntentConfig = { [key: string]: string }

type TConfig = {
  intents: Array<TIntentConfig>
  auth?: any
}
