import Command from '@oclif/command'
import { spawn } from 'child_process'
import * as path from 'path'

export default abstract class extends Command {
  runLegacy(cmdArgs: any[]) {
    const cliEntry = path.join(__dirname, '../node_modules/@bearer/bearer-cli/dist/bin/index.js')
    this.debug('Legacy command arguments', JSON.stringify(cmdArgs))
    spawn('node', [cliEntry, ...cmdArgs], { stdio: 'inherit' })
  }
}
