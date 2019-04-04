import BaseCommand from '../base-command'
import Cipher from '@bearer/security'

export default class Encrypt extends BaseCommand {
  static hidden = true
  static description = 'encrypt using bearer security'

  static flags = {
    ...BaseCommand.flags
  }

  static args = [{ name: 'encryptionKey', required: true }, { name: 'message', required: true }]

  async run() {
    const { args } = this.parse(Encrypt)
    const cipher = new Cipher(args.encryptionKey)
    this.log(cipher.encrypt(args.message))
  }
}
