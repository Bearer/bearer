import Command from '../BaseCommand'

type Constructor<T> = new (...args: any[]) => T
type TCommand = InstanceType<Constructor<Command>>

export function RequireScenarioFolder() {
  return function(_target: any, _propertyKey: string | symbol, descriptor: PropertyDescriptor) {
    let originalMethod = descriptor.value
    descriptor.value = async function(this: TCommand) {
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
    descriptor.value = async function(this: TCommand) {
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
  return function(_target: any, _propertyKey: string | symbol, descriptor: PropertyDescriptor) {
    let originalMethod = descriptor.value
    descriptor.value = async function(this: TCommand) {
      const { authorization, ExpiresAt } = this.bearerConfig.bearerConfig

      if (authorization && authorization.AuthenticationResult) {
        try {
          if (ExpiresAt < Date.now()) {
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
      } else {
        const error =
          this.colors.bold('⚠️ It looks like you are not logged in\n') +
          this.colors.yellow(this.colors.italic('Please run: bearer login'))
        this.error(error)
        return descriptor
      }
    }
    return descriptor
  }
}

async function refreshMyToken(command: TCommand): Promise<boolean | Error> {
  const { RefreshToken } = command.bearerConfig.bearerConfig.authorization.AuthenticationResult!

  // try {
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
      return true
    }
    case 401: {
      throw new UnauthorizedRefreshTokenError()
    }
    default: {
      throw new UnexpectedRefreshTokenError(JSON.stringify(body, undefined, 2))
    }
  }
}

class UnauthorizedRefreshTokenError extends Error {
  constructor() {
    super('Something went wrong, please run bearer login and try again')
  }
}

class UnexpectedRefreshTokenError extends Error {
  constructor(body: any) {
    super(`body : ${body}`)
  }
}
