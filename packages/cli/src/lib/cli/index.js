class CLI {
  constructor(program, emitter, config) {
    this.program = program
    this.emitter = emitter
    this.config = config
  }

  parse(argv) {
    this.program.parse(argv)

    // Display help if no command specified
    if (!argv.slice(2).length) {
      this.program.outputHelp()
      return
    }
  }

  use(cliCommand) {
    cliCommand.useWith(this.program, this.emitter, this.config)
  }
}

module.exports = { CLI }
