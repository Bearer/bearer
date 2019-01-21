import * as archiver from 'archiver'
import * as fs from 'fs-extra'
import * as path from 'path'

import BaseCommand from '../../base-command'
import { RequireScenarioFolder } from '../../utils/decorators'
import prepareConfig from '../../utils/prepare-config'

const CONFIG_FILE = 'bearer.config.json'
const HANDLER_NAME = 'index'
const HANDLER_NAME_WITH_EXT = [HANDLER_NAME, 'js'].join('.')

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
      const intents = await this.retrieveIntents()

      const { config, handlers } = buildLambdaDefinitions(intents)
      const finalConfig = await this.retrieveAuthConfig(config)

      archive.append(JSON.stringify(finalConfig, null, 2), { name: CONFIG_FILE })

      // add handlers
      this.debug(`Generated config: ${JSON.stringify(config, null, 2)}`)
      archive.append(handlers, { name: HANDLER_NAME_WITH_EXT })

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

  async retrieveAuthConfig(config: { intents: Record<string, string>[] }): Promise<TBearerConfig> {
    // generate config
    const content = await fs.readFile(this.locator.authConfigPath, { encoding: 'utf8' })
    return {
      ...config,
      auth: JSON.parse(content)
    }
  }

  // TODO: rewrite this using TS AST
  async retrieveIntents(): Promise<TIntentNames> {
    try {
      const config = await prepareConfig(
        this.locator.authConfigPath,
        this.bearerConfig.scenarioUuid,
        this.locator.srcIntentsDir
      )
      return config.intents
    } catch (e) {
      throw e
    }
  }
}

type TIntentNames = string[]

type TBearerConfig = {
  intents: Record<string, string>[]
  auth?: any
}

type TLambdaDefinition = {
  config: {
    intents: Record<string, string>[]
  }
  handlers: string
}

export function buildLambdaDefinitions(intents: TIntentNames): TLambdaDefinition {
  return {
    config: {
      intents: intents.reduce(
        (acc, intentName) => {
          acc.push({ [intentName]: [HANDLER_NAME, intentName].join('.') })
          return acc
        },
        [] as Record<string, string>[]
      )
    },
    handlers: buildLambdaIndex(intents)
  }
}

function buildLambdaIndex(intents: TIntentNames): string {
  return intents
    .map((intent, index) => {
      const intentConstName = `intent${index}`
      const intentInstance = `${intentConstName}Instance`
      return `const ${intentConstName} = require("./dist/${intent}").default;
const ${intentInstance} = new ${intentConstName}()
module.exports['${intent}'] = ${intentInstance}.intentType.intent(${intentInstance}.action);
`
    })
    .join('\n')
}
