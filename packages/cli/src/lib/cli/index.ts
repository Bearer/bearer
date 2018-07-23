import Locator from '../locationProvider'

export class CLI {
  constructor(
    private readonly program,
    private readonly emitter,
    private readonly config
  ) {}

  parse(argv) {
    this.program.parse(argv)

    // Display help if no command specified
    if (!argv.slice(2).length) {
      this.program.outputHelp()
      return
    }
  }

  use(cliCommand) {
    cliCommand.useWith(
      this.program,
      this.emitter,
      this.config,
      new Locator(this.config)
    )
  }
}
