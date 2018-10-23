class BearerContext {
  private state: { [key: string]: any } = {}
  private subscribers: Array<any> = []

  private _setupId: string

  get setupId(): string {
    return this._setupId
  }

  set setupId(setupId) {
    this._setupId = setupId
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
