import * as archiver from 'archiver'
import axios from 'axios'
import * as fs from 'fs-extra'
import * as globby from 'globby'
import * as Listr from 'listr'

import BaseCommand from '../base-command'
import { ensureFreshToken, RequireLinkedIntegration, RequireIntegrationFolder } from '../utils/decorators'
import { ensureFolderExists } from '../utils/helpers'
export default class Push extends BaseCommand {
  static description = 'Deploy Integration to Bearer Platform'

  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  @RequireIntegrationFolder()
  @RequireLinkedIntegration()
  @ensureFreshToken()
  async run() {
    ensureFolderExists(this.locator.buildArtifactDir, true)
    const archivePath = this.locator.buildArtifactResourcePath('integration.zip')
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
      this.success(`🐻 Integration successfully pushed.\n`)
      this.log(
        // tslint:disable-next-line:prefer-template
        `Your Integration will be available shortly here: ` +
          this.colors.bold(
            `${this.bearerConfig.DeveloperPortalUrl}integrations/${this.bearerConfig.integrationUuid}/preview`
          )
      )
      this.log(
        // tslint:disable-next-line:prefer-template
        `\nIn the meantime you can follow the deployment here: ` +
          this.colors.bold(
            `${this.bearerConfig.DeveloperPortalUrl}integrations/${this.bearerConfig.integrationUuid}/logs`
          )
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
        'functions/**/*.ts',
        'functions/tsconfig.json',
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
    return this.integrationClient.getIntegrationArchiveUploadUrl(
      this.bearerConfig.orgId!,
      this.bearerConfig.integrationId!
    )
  }

  async transfer(archivePath: string): Promise<boolean> {
    try {
      const url = await this.getSignedUrl()
      this.debug(url)
      const file = fs.readFileSync(archivePath)
      await axios.put(url, file, { headers: { 'Content-Type': 'application/zip' } })
      return true
    } catch (e) {
      if (e.response) {
        this.debug(e.response)
        switch (e.response.status) {
          case 401: {
            this.error(
              `Unauthorized to push, please visit ${this.bearerConfig.DeveloperPortalUrl}integrations/${
                this.bearerConfig.integrationUuid
              } to confirm you have access to ${this.colors.bold(this.bearerConfig.integrationUuid)} integration.`
            )
          }
          default: {
            this.log(e.response.data)
          }
        }
      }
      this.error(e)
      return false
    }
  }
}
