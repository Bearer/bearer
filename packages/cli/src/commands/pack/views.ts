import { flags } from '@oclif/command'
import { spawn } from 'child_process'
import * as fs from 'fs-extra'
import * as path from 'path'

import BaseCommand from '../../BaseCommand'
import { ScenarioBuildEnv } from '../../types'
import { RequireScenarioFolder } from '../../utils/decorators'

const scenarioId = 'scenario-id'
const scenarioUuid = 'scenario-uuid'
const cdnHost = 'cdn-host'

export default class PackViews extends BaseCommand {
  static description = 'Pack scenario views'
  static hidden = true
  static flags = {
    ...BaseCommand.flags,
    [scenarioUuid]: flags.string({
      required: true,
      description: 'Scenario unique identifier'
    }),
    [scenarioId]: flags.string({
      required: true,
      description: 'stencil scenario namespace'
    }),
    [cdnHost]: flags.string({
      required: true,
      description: 'Host url where views are uploade to (ex: https:static.bearer.sh/123456/attach-pull/dist/78901/'
    })
  }

  @RequireScenarioFolder()
  async run() {
    const { flags, args } = this.parse(PackViews)

    const config = this.bearerConfig
    const env: ScenarioBuildEnv = {
      ...process.env,
      BEARER_SCENARIO_ID: flags[scenarioUuid],
      BEARER_SCENARIO_TAG_NAME: flags[scenarioId],
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

  async buildStencil(env: ScenarioBuildEnv) {
    return new Promise((resolve, reject) => {
      const build = spawn('yarn', ['stencil', 'build'], { env, cwd: this.locator.buildViewsDir })

      build.stdout.on('data', data => {
        this.debug(`build scenario => stdout: ${data}`)
      })

      build.stderr.on('data', data => {
        this.debug(`build scenario => stderr: ${data}`)
      })

      build.on('close', code => {
        if (code === 0) {
          resolve()
        } else {
          reject(new Error("Can't build scenario views. please check logs"))
        }
      })
    })
  }
}
