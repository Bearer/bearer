const EventEmitter = require('events')
const term = require('terminal-kit').terminal

class Emitter {
  constructor() {
    this.emitter = new EventEmitter()
  }

  emit(name, args) {
    if (process.env.BEARER_DEBUG === '*') {
      term.white('Bearer event: ')
      term.yellow(name)
      term(' ')
      term.green(JSON.stringify(args))
      term('\n')
    }
    this.emitter.emit(name, args)
  }

  on(...args) {
    this.emitter.on(...args)
  }
}

module.exports = Emitter
