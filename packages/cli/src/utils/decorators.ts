import Command from '../base-command'

type Constructor<T> = new (...args: any[]) => T
type TCommand = InstanceType<Constructor<Command>>

export function skipIfNoViews() {
  return function(_target: any, _propertyKey: string | symbol, descriptor: PropertyDescriptor) {
    const originalMethod = descriptor.value
    descriptor.value = async function(this: TCommand) {
      if (this.hasViews) {
        await originalMethod.apply(this, arguments)
      } else {
        this.debug('No views present, skipping')
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
      if (this.bearerConfig.isIntegrationLocation) {
        await originalMethod.apply(this, arguments)
      } else {
        this.error('This command must be run within a integration folder.')
      }
    }
    return descriptor
  }
}

// tslint:disable-next-line:function-name
export function RequireLinkedIntegration() {
  return function(_target: any, _propertyKey: string | symbol, descriptor: PropertyDescriptor) {
    const originalMethod = descriptor.value
    descriptor.value = async function(this: TCommand) {
      if (this.bearerConfig.hasIntegrationLinked) {
        await originalMethod.apply(this, arguments)
      } else {
        const error =
          this.colors.bold('You integration must be linked before running this command\n') +
          this.colors.yellow(this.colors.italic('Please run: bearer link'))
        this.error(error)
      }
    }
    return descriptor
  }
}

export function ensureFreshToken() {
  return function(_target: any, _propertyKey: string | symbol, descriptor: PropertyDescriptor) {
    const originalMethod = descriptor.value
    descriptor.value = async function(this: TCommand) {
      const { expires_at } = (await this.bearerConfig.getToken()) || { expires_at: null }

      if (expires_at) {
        try {
          if (expires_at < Date.now()) {
            this.ux.action.start('Refreshing token')
            await refreshMyToken(this)
            this.ux.action.stop()
          }
        } catch (error) {
          this.ux.action.stop(`Failed`)
          this.error(error.message)
        }
        // TS is complaining with TS2445, commenting out this for now
        // this.debug('Running original method')
        await originalMethod.apply(this, arguments)
        return descriptor
      }

      const error =
        this.colors.bold('⚠️ It looks like you are not logged in\n') +
        this.colors.yellow(this.colors.italic('Please run: bearer login'))
      this.error(error)
      return descriptor
    }
    return descriptor
  }
}

async function refreshMyToken(command: TCommand): Promise<boolean | Error> {
  // TODO: rework refresh mechanism
  return true
}
