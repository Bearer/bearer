export interface IBearerConfig {
  integrationHost?: string
  clientId?: string
  loadingComponent?: string
}

export default class BearerConfig {
  integrationHost: string = 'BEARER_INTEGRATION_HOST'
  authorizationHost: string = 'BEARER_AUTHORIZATION_HOST'
  clientId: string = ''
  loadingComponent: string
  postRobotLogLevel: 'debug'| 'info'| 'warn'| 'error'  = 'error' 

  constructor(options: IBearerConfig = {}) {
    this.update(options)
  }

  update(options: IBearerConfig) {
    Object.keys(options).forEach(key => {
      this[key] = options[key]
    })
  }
}
