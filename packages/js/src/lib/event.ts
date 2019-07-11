export class EventEmitter {
  private readonly events: Record<string, Function[]> = {}

  on = (event: string, listener: Function) => {
    if (typeof this.events[event] !== 'object') {
      this.events[event] = []
    }

    this.events[event].push(listener)
  }

  clearListeners = (event: string) => {
    this.events[event] = []
  }

  removeListener = (event: string, listener: Function) => {
    if (this.events[event]) {
      const idx = this.events[event].indexOf(listener)
      if (idx > -1) {
        this.events[event].splice(idx, 1)
      }
    }
  }

  emit = (event: string, ...args: any[]) => {
    if (this.events[event]) {
      const listeners = this.events[event].slice()
      const length = listeners.length

      for (let i = 0; i < length; i++) {
        listeners[i].apply(this, args)
      }
    }
  }

  once = (event: string, listener: Function) => {
    const g = (...args: any[]) => {
      this.removeListener(event, g)
      listener.apply(this, args)
    }
    this.on(event, g)
  }
}
