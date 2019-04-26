import BaseCommand from '../base-command'
import Login from '../actions/login'

export default class LoginCommand extends BaseCommand {
  static description = 'login using Bearer credentials'

  static flags = {
    ...BaseCommand.flags
  }
  async run() {
    await Login(this)
  }
}
