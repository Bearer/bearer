import BaseCommand from '../../base-command'
import { RequireIntegrationFolder } from '../../utils/decorators'

import GenerateComponent from './component'
import GenerateFunction from './function'

const enum TType {
  Component,
  Function
}

export default class GenerateIndex extends BaseCommand {
  static description = 'generate Function or Component'
  static aliases = ['g']

  static flags = { ...BaseCommand.flags }

  static args = []

  @RequireIntegrationFolder()
  async run() {
    const { flags } = this.parse(GenerateIndex)
    const pathParams = flags.path ? ['--path', flags.path] : []
    const type = this.hasViews ? await this.askForType() : TType.Function

    switch (type) {
      case TType.Function: {
        return GenerateFunction.run([...pathParams])
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
        choices: [{ value: TType.Function, name: 'Function' }, { value: TType.Component, name: 'Component' }]
      }
    ])
    return type
  }
}
