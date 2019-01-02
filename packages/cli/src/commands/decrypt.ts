import BaseCommand from '../base-command'
import Cipher from '@bearer/security'

export default class Decrypt extends BaseCommand {
  static description = 'Decrypt using bearer security'
  static hidden = true

  static flags = {
    ...BaseCommand.flags
  }

  static args = [{ name: 'secretKey', required: true }, { name: 'decryptedMessage', required: true }]

  async run() {
    const { args } = this.parse(Decrypt)
    const cipher = new Cipher(args.secretKey)
    this.log(cipher.decrypt(args.decryptedMessage))
  }
}
