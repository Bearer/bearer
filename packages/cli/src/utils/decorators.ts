import * as inquirer from 'inquirer'
import Command from '../base-command'

import Link from '../actions/link'
import Create from '../actions/createIntegration'

type Constructor<T> = new (...args: any[]) => T
type TCommand = InstanceType<Constructor<Command>>

export function skipIfNoViews(displayError = false) {
  return function(_target: any, _propertyKey: string | symbol, descriptor: PropertyDescriptor) {
    const originalMethod = descriptor.value
    descriptor.value = async function(this: TCommand) {
      if (this.hasViews) {
        await originalMethod.apply(this, arguments)
      } else {
        if (displayError) {
          this.error(
            // tslint:disable-next-line
            'This integration does not contain any views. If you want to use this command please generate a new integration with --withViews flag'
          )
        } else {
          this.debug('No views present, skipping')
        }
      }
    }
    return descriptor
  }
}

// tslint:disable-next-line:function-name
export function RequireIntegrationFolder() {
  return function(_target: any, _propertyKey: string | symbol, descriptor: PropertyDescriptor) {
    const originalMethod = descriptor.value
    descriptor.value = async function(this: TCommand) {
      if (this.isIntegrationLocation) {
        await originalMethod.apply(this, arguments)
      } else {
        this.warn(
          // tslint:disable-next-line: max-line-length
          `We couldn't find any auth.config.json file, please make sure this file exists at the root of your integration`
        )
        this.error('This command must be run within a integration folder.')
      }
    }
    return descriptor
  }
}

// tslint:disable-next-line:function-name
export function RequireLinkedIntegration(prompt = true) {
  return function(_target: any, _propertyKey: string | symbol, descriptor: PropertyDescriptor) {
    const originalMethod = descriptor.value
    descriptor.value = async function(this: TCommand) {
      if (!this.bearerConfig.hasIntegrationLinked) {
        if (!prompt) {
          this.error('Can not run this command, please run link command before')
        }
        const { choice } = await inquirer.prompt([
          {
            name: 'choice',
            message: "Your integration isn't linked, what would you like to do?",
            type: 'list',
            choices: [
              {
                name: 'Create a new integration',
                value: 'create'
              },
              {
                name: 'Select an integration from my list',
                value: 'select'
              }
            ]
          }
        ])
        switch (choice) {
          case 'create':
            await Create(this, { link: true })
            break
          case 'select':
            await Link(this)
          default:
            break
        }
      }

      await originalMethod.apply(this, arguments)
    }
    return descriptor
  }
}
