class BearerContext {
  private state: { [key: string]: any } = {}
  private subscribers: Array<any> = []

  constructor() {
    console.log('[BEARER]', 'BearerContext init')
  }

  private _setupId: string
  private _configId: string

  get setupId(): string {
    return this._setupId
  }

  set setupId(setupId) {
    console.log('[BEARER]', 'setSetupId', setupId)
    this._setupId = setupId
  }

  get configId(): string {
    return this._configId
  }

  set configId(configId) {
    console.log('[BEARER]', 'setConfigId', configId)
    this._configId = configId
  }

  subscribe = (component: any) => {
    this.subscribers.push(component)
  }

  unsubscribe = (component: any) => {
    this.subscribers.filter(subscriber => subscriber === component)
  }

  update = (field, value) => {
    this.state[field] = value
    this.subscribers.map(component => {
      component.bearerUpdateFromState(this.state)
    })
  }
}

export default new BearerContext()
