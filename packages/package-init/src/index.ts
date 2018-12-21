import { exec } from 'child_process'
import { cli } from 'cli-ux'
import * as json from 'comment-json'
import * as fs from 'fs'
import { get, set } from 'lodash'

import * as path from 'path'
import * as util from 'util'

const spawnPromise = util.promisify(exec)

import { Command, flags } from '@oclif/command'
// import * as inquirer from 'inquirer'

// const enum TOptions {
//   typescript = 'typescript',
//   commitizen = 'commitizen',
//   tslint = 'tslint',
//   prettier = 'prettier',
//   jest = 'jest'
// }

class BearerPackageInit extends Command {
  static description = 'describe the command here'

  static flags = {
    path: flags.string({ char: 'p', default: process.cwd() }),
    force: flags.boolean({ char: 'f', default: false })
  }

  static args = [{ name: 'file' }]
  private cwd!: string
  private force!: boolean

  async run() {
    const { flags } = this.parse(BearerPackageInit)
    this.cwd = path.resolve(flags.path!)
    this.force = flags.force
    // add typescript
    // add husky
    // add commitizen
    // add lint-staged => prettier
    // const { choices } = await inquirer.prompt<{ choices: TOptions[] }>([
    //   {
    //     type: 'checkbox',
    //     choices: [
    //       {
    //         name: 'Typescript',
    //         value: TOptions.typescript,
    //         checked: true
    //       },
    //       {
    //         name: 'Commitizen (convential commits)',
    //         value: TOptions.commitizen,
    //         checked: true
    //       },
    //       {
    //         name: 'tslint (Typescript linter)',
    //         value: TOptions.tslint,
    //         checked: true
    //       },
    //       {
    //         name: 'prettier (format your code)',
    //         value: TOptions.prettier,
    //         checked: true
    //       },
    //       {
    //         name: 'jest (test your code)',
    //         value: TOptions.jest,
    //         checked: true
    //       }
    //     ],
    //     message: `Select things to add`,
    //     name: 'choices'
    //   }
    // ])
    // const options = new Set(choices)
    // add dependencies
    const dependencies = [
      'typescript',
      'husky',
      'lint-staged',
      'commitlint',
      '@commitlint/config-conventional',
      '@commitlint/cli',
      'cz-conventional-changelog',
      'tslint',
      BEARER_TSLINT,
      BEARER_TSCONFIG,
      'tslint-config-prettier',
      'prettier',
      'jest',
      'ts-jest',
      '@types/jest'
    ]
    try {
      await this.installDependencies(dependencies)
      await this.initTypescriptProject()
      await this.addHooksAndScripts()
      await this.addTsLintConfig()
      await this.setupJest()
    } catch (e) {
      cli.action.stop()
      return this.error(e)
    }
    // update package json
    // update ts config
  }

  async installDependencies(dependencies: string[]) {
    return this.withLoader(`Installing depenencies:\n  * ${dependencies.join('\n  * ')}`, () =>
      this.runCommand(`yarn add -D ${dependencies.join(' ')}`)
    )
  }

  async addHooksAndScripts() {
    const packageFile = path.join(this.cwd, 'package.json')
    try {
      cli.action.start('Adding hooks')
      const projectPackage: TPackage = JSON.parse(fs.readFileSync(packageFile, { encoding: 'utf8' }))
      if (!get(projectPackage, LINT_STAGED_KEY)) {
        set(projectPackage, LINT_STAGED_KEY, LINT_STAGED)
      } else {
        this.log('lint-staged already setup, skipping')
      }

      if (!get(projectPackage, COMMIT_LINT_CONFIG_KEY)) {
        set(projectPackage, COMMIT_LINT_CONFIG_KEY, COMMIT_LINT_CONFIG)
      } else {
        this.log('lint-staged already setup, skipping')
      }
      if (!get(projectPackage, LINT_STAGED_HOOK)) {
        set(projectPackage, LINT_STAGED_HOOK, LINT_STAGED_HOOK_VALUE)
      } else {
        this.log('lint-staged hook already setup, skipping')
      }
      if (!get(projectPackage, COMMIT_LINT_HOOK)) {
        set(projectPackage, COMMIT_LINT_HOOK, COMMIT_LINT_HOOK_VALUE)
      } else {
        this.log('commitlint hook already setup, skipping')
      }
      set(projectPackage, 'scripts.start', 'tsc --watch')
      set(projectPackage, 'scripts.build', 'tsc -p tsconfig.json')
      set(projectPackage, 'scripts.clean', 'rm -rf lib')
      set(projectPackage, 'scripts.prepack', 'yarn clean && yarn build')
      set(projectPackage, 'scripts.test', 'jest')
      set(projectPackage, 'scripts.test:ci', 'jest --coverage')
      fs.writeFileSync(packageFile, JSON.stringify(projectPackage, null, 2))
      fs.writeFileSync(
        path.join(this.cwd, 'commitlint.config.js'),
        "module.exports = { extends: ['@commitlint/config-conventional'] }"
      )

      cli.action.stop()
    } catch (e) {
      throw e
    }
  }

  async initTypescriptProject() {
    return this.withLoader('Init Typescript stuff', async () => {
      await this.runCommand('yarn tsc --init')
      const configFile = path.join(this.cwd, 'tsconfig.json')
      const src = path.join(this.cwd, 'src')
      const config = json.parse(fs.readFileSync(configFile, { encoding: 'utf8' }), undefined, true)

      set(config, 'extends', BEARER_TSCONFIG)
      fs.writeFileSync(configFile, json.stringify(config, null, 2))
      if (!fs.existsSync(src)) {
        fs.mkdirSync(src)
        fs.writeFileSync(path.join(src, 'index.ts'), '')
      }
    })
  }

  async addTsLintConfig() {
    return this.withLoader('Init TSlint stuff', async () => {
      const configFile = path.join(this.cwd, 'tslint.json')
      if (!fs.existsSync(configFile)) {
        await this.runCommand('yarn tslint --init')
      }
      const config = json.parse(fs.readFileSync(configFile, { encoding: 'utf8' }), undefined, true)
      // as soon a we have created a bearer tslint we use it herer
      set(config, 'extends', BEARER_TSLINT)
      fs.writeFileSync(configFile, json.stringify(config, null, 2))
      fs.writeFileSync(path.join(this.cwd, '.prettierrc'), json.stringify(prettierConfig, null, 2))
    })
  }

  async setupJest() {
    // check if a jest config is already present in package.json
    if (!fs.existsSync(path.join(this.cwd, 'jest.config.js'))) {
      return this.withLoader('Init Jest stuff', async () => {
        await this.runCommand('yarn ts-jest config:init')
      })
    }
  }

  async withLoader(title: string, block: () => Promise<any>) {
    try {
      cli.action.start(title)
      await block()
      cli.action.stop()
    } catch (e) {
      if (this.force) {
        console.error('Error on ', title, '\n', e.toString())
      } else {
        throw e
      }
    }
  }

  async runCommand(command: string) {
    await spawnPromise(command, {
      cwd: this.cwd
    })
  }
}

type TPackage = {
  config?: {
    commitizen?: {
      path: string
    }
  }
  ['lint-staged']?: {
    [key: string]: string[]
  }
}

const COMMIT_LINT_HOOK = 'husky.hooks.commit-msg'
const COMMIT_LINT_HOOK_VALUE = 'commitlint -E HUSKY_GIT_PARAMS'
const COMMIT_LINT_CONFIG = './node_modules/cz-conventional-changelog'
const COMMIT_LINT_CONFIG_KEY = 'config.commitizen.path'
const LINT_STAGED_HOOK = 'husky.hooks.pre-commit'
const LINT_STAGED_HOOK_VALUE = 'lint-staged'
const LINT_STAGED_KEY = 'lint-staged'
const LINT_STAGED = {
  '*.{css,md,tsx,ts}': ['prettier --write', 'tslint -c tslint.json --fix', 'git add']
}
const BEARER_TSLINT = '@bearer/tslint-config'
const BEARER_TSCONFIG = '@bearer/tsconfig'
const prettierConfig = {
  semi: false,
  singleQuote: true,
  printWidth: 120,
  tabWidth: 2
}

export = BearerPackageInit
