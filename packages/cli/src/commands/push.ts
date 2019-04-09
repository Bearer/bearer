import * as archiver from 'archiver'
import axios from 'axios'
import * as fs from 'fs-extra'
import * as globby from 'globby'
import * as Listr from 'listr'
import serviceClient from '@bearer/bearer-cli/lib/lib/serviceClient'

import BaseCommand from '../base-command'
import { ensureFreshToken, RequireLinkedIntegration, RequireIntegrationFolder } from '../utils/decorators'
import { ensureFolderExists } from '../utils/helpers'
import { IntegrationClient } from '../utils/integration-client'
import { AUTH_CONFIG_FILENAME } from '../utils/locator'

export default class Push extends BaseCommand {
  static description = 'deploy Integration to Bearer'
  private integrationClient!: IntegrationClient
  // TODO: fix typing
  private serviceClient!: any

  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  @RequireIntegrationFolder()
  @ensureFreshToken()
  @RequireLinkedIntegration()
  async run() {
    const { Username, Password } = await this.fetchLoginInformation()
    this.serviceClient = serviceClient(this.constants.IntegrationServiceUrl)
    const data = await this.serviceClient.login({ Username, Password })
    this.debug('auth info : %j', data)

    this.integrationClient = new IntegrationClient(
      this.constants.DeploymentUrl,
      data.body.authorization.AuthenticationResult.IdToken,
      this.config.version
    )

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

  async getSignedUrl(): Promise<string> {
    return this.integrationClient.getIntegrationArchiveUploadUrl(this.bearerConfig.bearerUid)
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

  fetchLoginInformation = async () => {
    const { BEARER_TOKEN, BEARER_EMAIL } = process.env
    if (BEARER_TOKEN && BEARER_EMAIL) {
      return {
        Username: BEARER_EMAIL,
        Password: BEARER_TOKEN
      }
    }
    const { data } = await this.devPortalClient.request<{
      currentUser: { email: string; infrastructure: { password: string } }
    }>({
      query: QUERY
    })
    if (data.data) {
      return { Username: data.data.currentUser.email, Password: data.data.currentUser.infrastructure.password }
    }
    throw 'Fetch credentials error'
  }
}

const QUERY = `
query CLIPush {
  currentUser {
    email
    infrastructure {
      password
    }
  }
}
`
