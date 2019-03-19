import { EventEmitter } from 'events'

const term = require('terminal-kit').terminal

class Emitter {
  private emitter: EventEmitter
  constructor() {
    this.emitter = new EventEmitter()
  }

  emit(name, args) {
    if (/bearer:/.test(process.env.DEBUG || '')) {
      term.white('Bearer event: ')
      term.yellow(name)
      term(' ')
      term.green(JSON.stringify(args))
      term('\n')
    }
    this.emitter.emit(name, args)
  }

  on(event: string | symbol, listener: (...args: any[]) => void) {
    this.emitter.on(event, listener)
  }
}

export default Emitter
