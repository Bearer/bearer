import Command from '../BaseCommand'
interface ICommand extends Command {}

export function RequireScenarioFolder() {
  return function(_target: any, _propertyKey: string | symbol, descriptor: PropertyDescriptor) {
    let originalMethod = descriptor.value
    descriptor.value = function(this: ICommand) {
      if (this.bearerConfig.isScenarioLocation) {
        originalMethod.apply(this, arguments)
      } else {
        this.error('This command must be run within a scenario folder.')
      }
    }
    return descriptor
  }
}
