import * as archiver from 'archiver'
import * as fs from 'fs-extra'
import * as path from 'path'

import BaseCommand from '../../base-command'
import { RequireIntegrationFolder } from '../../utils/decorators'
import prepareConfig from '../../utils/prepare-config'

const CONFIG_FILE = 'bearer.config.json'
const HANDLER_NAME = 'index'
const HANDLER_NAME_WITH_EXT = [HANDLER_NAME, 'js'].join('.')

export default class PackFunctions extends BaseCommand {
  static description = 'Pack integration functions'
  static hidden = true
  static flags = {
    ...BaseCommand.flags
  }

  static args = [{ name: 'ARCHIVE_PATH', required: true }]

  @RequireIntegrationFolder()
  async run() {
    // we are assuming prepare and build have been run previously
    const { args } = this.parse(PackFunctions)
    const { ARCHIVE_PATH } = args
    const target = path.resolve(ARCHIVE_PATH)
    this.debug(`Zipping to: ${target}`)

    const output = fs.createWriteStream(target)
    const archive = archiver('zip', {
      zlib: { level: 9 }
    })
    const archiveEntries: archiver.EntryData[] = []
    archive.pipe(output)

    const archived = new Promise<archiver.EntryData[]>((resolve, reject) => {
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
      // add functions
      archive.directory(this.locator.buildFunctionsDir, false)
      // add CONFIG
      const functions = await this.retrieveFunctions()

      const { config, handlers } = buildLambdaDefinitions(functions)
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

  async retrieveAuthConfig(config: { functions: Record<string, string>[] }): Promise<TBearerConfig> {
    // generate config
    const content = await fs.readFile(this.locator.authConfigPath, { encoding: 'utf8' })
    return {
      ...config,
      auth: JSON.parse(content)
    }
  }

  // TODO: rewrite this using TS AST
  async retrieveFunctions(): Promise<TFunctionNames> {
    try {
      const config = await prepareConfig(
        this.locator.authConfigPath,
        this.bearerConfig.integrationUuid,
        this.locator.srcFunctionsDir
      )
      return config.functions
    } catch (e) {
      throw e
    }
  }
}

type TFunctionNames = string[]

type TBearerConfig = {
  functions: Record<string, string>[]
  auth?: any
}

type TLambdaDefinition = {
  config: {
    functions: Record<string, string>[]
  }
  handlers: string
}

export function buildLambdaDefinitions(functions: TFunctionNames): TLambdaDefinition {
  return {
    config: {
      functions: functions.reduce(
        (acc, functionName) => {
          acc.push({ [functionName]: [HANDLER_NAME, functionName].join('.') })
          return acc
        },
        [] as Record<string, string>[]
      )
    },
    handlers: buildLambdaIndex(functions)
  }
}

function buildLambdaIndex(functions: TFunctionNames): string {
  return functions
    .map((func, index) => {
      const funcConstName = `func${index}`
      return `const ${funcConstName} = require("./dist/${func}").default;
module.exports['${func}'] = ${funcConstName}.init();
`
    })
    .join('\n')
}
