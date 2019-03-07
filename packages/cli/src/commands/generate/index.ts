import BaseCommand from '../../base-command'
import { RequireIntegrationFolder } from '../../utils/decorators'

import GenerateComponent from './component'
import GenerateIntent from './intent'

const enum TType {
  Component,
  Intent
}

export default class GenerateIndex extends BaseCommand {
  static description = 'Generate Intent or Component'
  static aliases = ['g']

  static flags = { ...BaseCommand.flags }

  static args = []

  @RequireIntegrationFolder()
  async run() {
    const { flags } = this.parse(GenerateIndex)
    const pathParams = flags.path ? ['--path', flags.path] : []
    const type = await this.askForType()

    switch (type) {
      case TType.Intent: {
        return GenerateIntent.run([...pathParams])
      }
      case TType.Component: {
        return GenerateComponent.run([...pathParams])
      }
    }
  }

  async askForType(): Promise<TType> {
    const { type } = await this.inquirer.prompt<{ type: TType }>([
      {
        message: 'What would you like to generate:',
        type: 'list',
        name: 'type',
        choices: [{ value: TType.Intent, name: 'Intent' }, { value: TType.Component, name: 'Component' }]
      }
    ])
    return type
  }
}
