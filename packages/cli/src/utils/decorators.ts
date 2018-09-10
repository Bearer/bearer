import Command from '../BaseCommand'

type Constructor<T> = new (...args: any[]) => T;
type TCommand = InstanceType<Constructor<Command>>

export function RequireScenarioFolder() {
  return function (_target: any, _propertyKey: string | symbol, descriptor: PropertyDescriptor) {
    let originalMethod = descriptor.value
    descriptor.value = async function (this: TCommand) {
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
  return function (_target: any, _propertyKey: string | symbol, descriptor: PropertyDescriptor) {
    let originalMethod = descriptor.value
    descriptor.value = async function (this: TCommand) {
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

export function ensureFreshToken() {
  return function (_target: any, _propertyKey: string | symbol, descriptor: PropertyDescriptor) {
    let originalMethod = descriptor.value
    descriptor.value = async function (this: TCommand) {
      const { authorization, ExpiresAt } = this.bearerConfig.bearerConfig

      if (authorization && authorization.AuthenticationResult) {
        try {
          if (ExpiresAt < Date.now()) {
            await refreshMyToken(this)
          }
        } catch (e) {
          console.log('[BEARER]', 'e', e)
          return this.error(e)
        }
        // TS is complaining with TS2445, commenting out this for now
        // this.debug('Running original method')
        await originalMethod.apply(this, arguments)
      } else {
        const error =
          this.colors.bold('⚠️ It looks like you are not connected\n') +
          this.colors.yellow(this.colors.italic('Please run: bearer login'))
        this.error(error)
      }
    }
    return descriptor
  }
}

async function refreshMyToken(command: TCommand) {
  const { RefreshToken } = command.bearerConfig.bearerConfig.authorization.AuthenticationResult!

  try {
    command.ux.action.start('Refreshing token')
    const { statusCode, body } = await command.serviceClient.refresh({ RefreshToken })

    switch (statusCode) {
      case 200: {
        const {
          authorization: { AuthenticationResult }
        } = body

        await command.bearerConfig.storeBearerConfig({
          ExpiresAt: Date.now() + AuthenticationResult.ExpiresIn * 1000,
          authorization: {
            AuthenticationResult: {
              ...AuthenticationResult,
              RefreshToken
            }
          }
        })
        break
      }
      case 401: {
        throw new Error('Unauthorized, please run bearer login')
      }
      default: {
        command.log('Error: ', body)
        throw new Error('Unexpected error')
      }
    }
  } catch (error) {
    throw error
  }
  command.ux.action.stop()
}
