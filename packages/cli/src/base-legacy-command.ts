import Command from '@oclif/command'

export default abstract class extends Command {
  runLegacy(cmdArgs: any[]) {
    const cli = require('@bearer/bearer-cli/lib/bin/index').default
    this.debug('Legacy command arguments', JSON.stringify(cmdArgs))
    cli(['null', 'null'].concat(cmdArgs))
  }
}
