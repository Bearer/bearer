import BaseCommand from '../../BaseCommand'
import { RequireScenarioFolder } from '../../utils/decorators'

import GenerateComponent from './component'
import GenerateIntent from './intent'

const enum TType {
  Component,
  Intent
}

export default class GenerateIndex extends BaseCommand {
  static description = 'Generate Intent or Component'

  static flags = {
    ...BaseCommand.flags
  }

  static args = []

  @RequireScenarioFolder()
  async run() {
    const type = await this.askForType()

    switch (type) {
      case TType.Intent: {
        return GenerateIntent.run([])
      }
      case TType.Component: {
        return GenerateComponent.run([])
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
