import { exec } from 'child_process'
import { cli } from 'cli-ux'
import * as path from 'path'
import * as util from 'util'
const spawnPromise = util.promisify(exec)

import { Command, flags } from '@oclif/command'
import * as inquirer from 'inquirer'

const enum TOptions {
  typescript = 'typescript',
  commitizen = 'commitizen',
  tslint = 'tslint',
  prettier = 'prettier',
  jest = 'jest'
}

class BearerPackageInit extends Command {
  static description = 'describe the command here'

  static flags = {
    path: flags.string({ char: 'p', default: process.cwd() })
  }

  static args = [{ name: 'file' }]

  async run() {
    const { flags } = this.parse(BearerPackageInit)
    const cwd = path.resolve(flags.path!)
    // add typescript
    // add husky
    // add commitizen
    // add lint-staged => prettier
    const { choices } = await inquirer.prompt<{ choices: TOptions[] }>([
      {
        type: 'checkbox',
        choices: [
          {
            name: 'Typescript',
            value: TOptions.typescript,
            checked: true
          },
          {
            name: 'Commitizen (convential commits)',
            value: TOptions.commitizen,
            checked: true
          },
          {
            name: 'tslint (Typescript linter)',
            value: TOptions.tslint,
            checked: true
          },
          {
            name: 'prettier (format your code)',
            value: TOptions.prettier,
            checked: true
          },
          {
            name: 'jest (test your code)',
            value: TOptions.jest,
            checked: true
          }
        ],
        message: `Select things to add`,
        name: 'choices'
      }
    ])
    const options = new Set(choices)
    // add dependencies
    const dependencies = [
      'tsc',
      'husky',
      'lint-staged',
      'commitlint',
      '@commitlint/config-conventional',
      '@commitlint/cli',
      'cz-conventional-changelog',
      'tslint',
      'tslint-config-prettier',
      'prettier'
    ].join(' ')
    try {
      cli.action.start(`Installing depenencies : ${[...options].join(', ')}`)
      await spawnPromise(`yarn add -D ${dependencies}`, {
        cwd
      })
      cli.action.stop()
    } catch (e) {
      return this.error(e)
    }
    // cli.action.start('Updating files')

    // cli.action.stop()
    // update package json
    // update ts config
  }
}

export = BearerPackageInit
