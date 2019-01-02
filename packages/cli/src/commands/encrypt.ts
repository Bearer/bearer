import BaseCommand from '../base-command'
import Cipher from '@bearer/security'

export default class Encrypt extends BaseCommand {
  static description = 'Encrypt using bearer security'
  static hidden = true

  static flags = {
    ...BaseCommand.flags
  }

  static args = [{ name: 'secretKey', required: true }, { name: 'message', required: true }]

  async run() {
    const { args } = this.parse(Encrypt)
    const cipher = new Cipher(args.secretKey)
    this.log(cipher.encrypt(args.message))
  }
}
