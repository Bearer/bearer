import * as archiver from 'archiver'
import axios from 'axios'
import * as fs from 'fs-extra'
import * as globby from 'globby'
import * as Listr from 'listr'

import BaseCommand from '../base-command'
import { RequireLinkedIntegration, RequireIntegrationFolder } from '../utils/decorators'
import { ensureFolderExists } from '../utils/helpers'
import { AUTH_CONFIG_FILENAME } from '../utils/locator'

export default class Push extends BaseCommand {
  static description = 'deploy integration to Bearer'

  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  @RequireIntegrationFolder()
  @RequireLinkedIntegration()
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
          this.colors.bold(`${this.constants.DeveloperPortalUrl}integrations/${this.bearerConfig.bearerUid}`)
      )
      this.log(
        // tslint:disable-next-line:prefer-template
        `\nIn the meantime you can follow the deployment here: ` +
          this.colors.bold(`${this.constants.DeveloperPortalUrl}integrations/${this.bearerConfig.bearerUid}/logs`)
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
        AUTH_CONFIG_FILENAME
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

  async getSignedUrl() {
    try {
      const { data } = await this.devPortalClient.request<{ integration: { url: string } }>({
        query: QUERY,
        variables: { buid: this.bearerConfig.BUID }
      })

      if (!data.data) {
        this.debug('received: %j', data)
        this.error(new Error('Integration not found please make sure you have correctly linked your integration'))
      }
      return data!.data!.integration.url
    } catch (e) {
      this.error(e)
      throw e
    }
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
          case 403:
          case 401: {
            this.error(
              `Unauthorized to push, please visit ${this.constants.DeveloperPortalUrl}integrations/${
                this.bearerConfig.bearerUid
              } to confirm you have access to ${this.colors.bold(this.bearerConfig.bearerUid)} integration.`
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

const QUERY = `
query CLIGetIntegrationUploadUrl($buid: String!) {
  integration(buid: $buid) {
    url: archiveUploadUrl
  }
}
`
