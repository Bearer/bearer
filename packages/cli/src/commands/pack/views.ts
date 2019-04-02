import { flags } from '@oclif/command'
import { spawn } from 'child_process'

import BaseCommand from '../../base-command'
import { IntegrationBuildEnv } from '../../types'
import { RequireIntegrationFolder, skipIfNoViews } from '../../utils/decorators'

const integrationId = 'integration-id'
const integrationUuid = 'integration-uuid'
const cdnHost = 'cdn-host'

export default class PackViews extends BaseCommand {
  static description = 'Pack integration views'
  static hidden = true
  static flags = {
    ...BaseCommand.flags,
    [integrationUuid]: flags.string({
      required: true,
      description: 'Integration unique identifier'
    }),
    [integrationId]: flags.string({
      required: true,
      description: 'stencil integration namespace'
    }),
    [cdnHost]: flags.string({
      required: true,
      description: 'Host url where views are uploade to (ex: https:static.bearer.sh/123456/attach-pull/dist/78901/'
    })
  }

  @skipIfNoViews()
  @RequireIntegrationFolder()
  async run() {
    const { flags } = this.parse(PackViews)

    const config = this.constants
    const env: IntegrationBuildEnv = {
      ...process.env,
      BEARER_INTEGRATION_ID: flags[integrationUuid],
      BEARER_INTEGRATION_TAG_NAME: flags[integrationId],
      BEARER_INTEGRATION_HOST: config.IntegrationServiceHost,
      BEARER_AUTHORIZATION_HOST: config.IntegrationServiceHost,
      CDN_HOST: flags[cdnHost]
    }

    try {
      const buildDestination = await this.buildStencil(env)
      this.success(`Packed views : ${buildDestination}`)
    } catch (e) {
      this.error(e)
    }
  }

  async buildStencil(env: IntegrationBuildEnv) {
    return new Promise((resolve, reject) => {
      const build = spawn('yarn', ['stencil', 'build'], { env, cwd: this.locator.buildViewsDir })

      build.stdout.on('data', data => {
        this.debug(`build integration => stdout: ${data}`)
      })

      build.stderr.on('data', data => {
        this.debug(`build integration => stderr: ${data}`)
      })

      build.on('close', code => {
        if (code === 0) {
          resolve()
        } else {
          reject(new Error("Can't build integration views. please check logs"))
        }
      })
    })
  }
}
