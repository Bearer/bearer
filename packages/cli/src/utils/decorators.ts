import Command from '../BaseCommand'
interface ICommand extends Command {}

export function RequireScenarioFolder() {
  return function(_target: any, _propertyKey: string | symbol, descriptor: PropertyDescriptor) {
    let originalMethod = descriptor.value
    descriptor.value = async function(this: ICommand) {
      if (this.bearerConfig.isScenarioLocation) {
        await originalMethod.apply(this, arguments)
      } else {
        this.error('This command must be run within a scenario folder.')
      }
    }
    return descriptor
  }
}

export function RequireLinkedScenario() {
  return function(_target: any, _propertyKey: string | symbol, descriptor: PropertyDescriptor) {
    let originalMethod = descriptor.value
    descriptor.value = async function(this: ICommand) {
      if (this.bearerConfig.hasScenarioLinked) {
        await originalMethod.apply(this, arguments)
      } else {
        const error =
          this.colors.bold('You scenario must be linked before running this command\n') +
          this.colors.yellow(this.colors.italic('Please run: bearer link'))
        this.error(error)
      }
    }
    return descriptor
  }
}
