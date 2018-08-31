import transpiler from '@bearer/transpiler/lib/bin/bearer-tst.js'
import { flags } from '@oclif/command'
import * as chokidar from 'chokidar'
import * as fs from 'fs-extra'
import * as Listr from 'listr'
import * as path from 'path'

import BaseCommand from '../../BaseCommand'
import installDependencies from '../../tasks/installDependencies'
import runNpmCommand from '../../tasks/runNpmClientCommand'
import { ScenarioBuildEnv } from '../../types'
import { RequireScenarioFolder } from '../../utils/decorators'

const skipInstall = 'skip-install'

export default class BuildViews extends BaseCommand {
  static description = 'Build scenario views'
  static hidden = true
  static flags = {
    ...BaseCommand.flags,
    [skipInstall]: flags.boolean({})
  }

  static args = []

  @RequireScenarioFolder()
  async run() {
    const { flags } = this.parse(BuildViews)

    const config = this.bearerConfig
    const cdnHost = `${config.CdnHost}/${config.orgId}/${config.scenarioId}/dist/${config.scenarioId}/`
    // TODO: throw error if missing information
    const env: ScenarioBuildEnv = {
      BEARER_SCENARIO_ID: config.scenarioUuid,
      BEARER_SCENARIO_TAG_NAME: config.scenarioId!,
      BEARER_INTEGRATION_HOST: config.IntegrationServiceHost,
      BEARER_AUTHORIZATION_HOST: config.IntegrationServiceHost,
      CDN_HOST: cdnHost,
      ...process.env
    }

    const tasks: Array<Listr.ListrTask> = [
      {
        title: 'Transpile views',
        task: async (_ctx: any, _task: any) => this.transpile()
      },
      {
        title: 'Link non TS files',
        task: async (_ctx: any, _task: any) => {
          const watcher = await this.watchNonTsFiles(this.locator.srcViewsDir, this.locator.buildViewsComponentsDir)
          // TODO: update here when watch mode required
          // if(flags.noWatch) { // or similar
          watcher.close()
          //}
        }
      },
      runNpmCommand({ name: 'Build scenario views', cwd: this.locator.buildViewsDir, command: 'stencil build', env })
    ]
    if (!flags[skipInstall]) {
      tasks.unshift(installDependencies({ cwd: this.locator.scenarioRoot }))
    }

    try {
      await new Listr(tasks).run()
      this.success('Built views')
    } catch (e) {
      this.error(e)
    }
  }

  transpile = async () => {
    const prefix = ['bearer', this.bearerConfig.scenarioId].join('-')
    const suffix = this.bearerConfig.orgId
    try {
      transpiler(['--no-watcher', '--prefix-tag', prefix, '--suffix-tag', suffix, '--no-process'])
    } catch (e) {
      this.error(e)
    }
    await setTimeout(() => {}, 5000)
  }

  watchNonTsFiles = async (watchedPath: string, destPath: string): Promise<chokidar.FSWatcher> => {
    return new Promise<chokidar.FSWatcher>((resolve, _reject) => {
      const callback = (error: any) => {
        if (error) {
          this.error(error)
        }
      }

      const watcher = chokidar.watch(watchedPath + '/**', {
        ignored: /\.tsx?$/,
        persistent: true,
        followSymlinks: false
      })

      watcher.on('ready', () => {
        resolve(watcher)
      })

      watcher.on('all', (event, filePath) => {
        const relativePath = filePath.replace(watchedPath, '')
        const targetPath = path.join(destPath, relativePath)
        // Creating symlink
        if (event === 'add') {
          this.debug('creating symlink', filePath, targetPath)
          fs.ensureSymlink(filePath, targetPath, callback)
        }
        // // Deleting symlink
        if (event === 'unlink') {
          this.debug('deleting symlink')
          fs.unlink(targetPath, err => {
            if (err) throw err
            this.debug(targetPath + ' was deleted')
          })
        }
      })
    })
  }
}
