import axios from 'axios'

import Command from '../base-command'
import { LOGIN_CLIENT_ID } from './constants'

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
      const { expires_at, refresh_token } = (await this.bearerConfig.getToken()) || {
        expires_at: null,
        refresh_token: null
      }

      if (expires_at && refresh_token) {
        try {
          if (expires_at < Date.now()) {
            this.ux.action.start('Refreshing token')
            await refreshMyToken(this, refresh_token)
            this.ux.action.stop()
          }
        } catch (error) {
          this.ux.action.stop(`Failed`)
          this.error(error.message)
        }
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

// tslint:disable-next-line variable-name
async function refreshMyToken(command: TCommand, refresh_token: string): Promise<boolean | Error> {
  // TODO: rework refresh mechanism
  const response = await axios.post(`${command.constants.LoginDomain}/oauth/token`, {
    refresh_token,
    grant_type: 'refresh_token',
    client_id: LOGIN_CLIENT_ID
  })
  await command.bearerConfig.storeToken({ ...response.data, refresh_token })
  return true
}
