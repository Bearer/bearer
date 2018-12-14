import { BearerWindow } from '@bearer/types'
declare const window: BearerWindow

export interface IBearerConfig {
  integrationHost?: string
  loadingComponent?: string
}

export default class BearerConfig {
  integrationHost = 'BEARER_INTEGRATION_HOST'
  authorizationHost = 'BEARER_AUTHORIZATION_HOST'
  loadingComponent: string
  postRobotLogLevel: 'debug' | 'info' | 'warn' | 'error' = 'error'

  get clientId(): string {
    return window.bearer && window.bearer.clientId
  }

  constructor(options: IBearerConfig = {}) {
    this.update(options)
  }

  update(options: IBearerConfig) {
    Object.keys(options).forEach(key => {
      this[key] = options[key]
    })
  }
}
