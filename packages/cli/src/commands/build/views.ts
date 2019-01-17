import transpiler from '@bearer/transpiler/lib/bin/bearer-tst.js'
import { flags } from '@oclif/command'
import * as chokidar from 'chokidar'
import * as fs from 'fs-extra'
import * as Listr from 'listr'
import * as path from 'path'

import BaseCommand from '../../base-command'
import installDependencies from '../../tasks/install-dependencies'

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
    const tasks: Listr.ListrTask[] = [
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
          // }
        }
      }
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

  transpile = () => {
    const prefix = ['bearer', this.bearerConfig.orgId].join('-')
    const suffix = this.bearerConfig.scenarioId
    try {
      transpiler(['--no-watcher', '--prefix-tag', prefix, '--suffix-tag', suffix, '--no-process'])
    } catch (e) {
      this.error(e)
    }
  }

  watchNonTsFiles = async (watchedPath: string, destPath: string): Promise<chokidar.FSWatcher> => {
    return new Promise<chokidar.FSWatcher>((resolve, _reject) => {
      const callback = (error: any) => {
        if (error) {
          this.error(error)
        }
      }

      const watcher = chokidar.watch(`${watchedPath}/**`, {
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
            this.debug(`${targetPath} was deleted`)
          })
        }
      })
    })
  }
}
