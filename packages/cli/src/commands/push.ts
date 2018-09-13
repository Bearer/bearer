import * as archiver from 'archiver'
import axios from 'axios'
import * as fs from 'fs-extra'
import * as globby from 'globby'
import * as Listr from 'listr'

import BaseCommand from '../BaseCommand'
import { ensureFreshToken, RequireLinkedScenario, RequireScenarioFolder } from '../utils/decorators'
import { ensureFolderExists } from '../utils/helpers'
export default class Push extends BaseCommand {
  static description = 'Deploy Scenario to Bearer Platform'

  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  @RequireScenarioFolder()
  @RequireLinkedScenario()
  @ensureFreshToken()
  async run() {
    ensureFolderExists(this.locator.buildArtifactDir, true)
    const archivePath = this.locator.buildArtifactResourcePath('scenario.zip')
    const tasks = [
      {
        title: 'Generate bundle',
        task: async (_ctx: any) => this.archive(archivePath)
      },
      {
        title: 'Transfer bundle',
        task: async (_ctx: any) => this.transfer(archivePath)
      }
    ]

    try {
      await new Listr(tasks).run()
      this.success(`🐻 Scenario successfully pushed.\n`)
      this.log(
        `Your scenario will be available soon at this location: ` +
        this.colors.bold(`${this.bearerConfig.DeveloperPortalUrl}scenarios/${this.bearerConfig.scenarioUuid}/preview`)
      )
      this.log(
        `\nIn the mean time you can follow the deployment here: ` +
        this.colors.bold(`${this.bearerConfig.DeveloperPortalUrl}scenarios/${this.bearerConfig.scenarioUuid}/logs`)
      )
    } catch (e) {
      this.error(e)
    }
  }

  async archive(archivePath: string): Promise<string> {
    return new Promise<any>(async (resolve, reject) => {
      const output = fs.createWriteStream(archivePath)
      const archive = archiver('zip', {
        zlib: { level: 9 } // Sets the compression level.
      })
      // pipe archive data to the file
      archive.pipe(output)
      const files = await globby([
        'views/**/*',
        'intents/**/*.ts',
        'intents/tsconfig.json',
        'yarn.lock',
        'package-json.lock',
        'spec.ts',
        'package.json',
        'auth.config.json'
      ])
      this.debug('Files to upload', files.join('\n'))

      if (files.length >= 100) {
        return reject(new Error('Too many files to bundle. Please re-run this command this DEBUG=*'))
      }

      output.on('close', () => {
        this.debug(`Archive created: ${archive.pointer() / 1024} Kb / ${archivePath}`)
        resolve(archivePath)
      })
      archive.on('error', (err: any) => {
        reject(err)
      })

      archive.on('warning', (err: any) => {
        if (err.code === 'ENOENT') {
          reject(err)
        } else {
          this.debug(err)
        }
      })

      files.map(file => {
        archive.file(file, { name: file })
      })
      archive.finalize()
    })
  }

  async getSignedUrl(): Promise<string> {
    return this.scenarioClient.getScenarioArchiveUploadUrl(this.bearerConfig.orgId!, this.bearerConfig.scenarioId!)
  }

  async transfer(archivePath: string): Promise<boolean> {
    try {
      const url = await this.getSignedUrl()
      this.debug(url)
      const file = fs.readFileSync(archivePath)
      await axios.put(url, file, { headers: { 'Content-Type': 'application/zip' } })
      return true
    } catch (e) {
      if (e.response && e.response.status === 401) {
        this.error('Unauthorized')
        return false
      } else {
        throw e
      }
    }
  }
}
