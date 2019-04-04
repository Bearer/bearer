import Locator from '../locationProvider'
import { Config } from '../types'

export class CLI {
  constructor(readonly program, private readonly emitter, private readonly config: Config) {}

  parse(argv) {
    this.program.parse(argv)

    // Display help if no command specified
    if (!argv.slice(2).length) {
      this.program.outputHelp()
      return
    }
  }

  use(cliCommand) {
    cliCommand.useWith(this.program, this.emitter, this.config, new Locator())
  }
}
