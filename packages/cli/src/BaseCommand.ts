import Command from '@oclif/command'
import * as inquirer from 'inquirer'

export default abstract class extends Command {
  get inquirer() {
    return inquirer
  }
}
