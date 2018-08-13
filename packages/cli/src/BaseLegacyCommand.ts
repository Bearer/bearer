import cli from '@bearer/bearer-cli-testing-match-issue/dist/bin/index'
import Command from '@oclif/command'

export default abstract class extends Command {
  runLegacy(cmdArgs: any[]) {
    this.debug('Legacy command arguments', JSON.stringify(cmdArgs))
    cli(['null', 'null'].concat(cmdArgs))
  }
}
