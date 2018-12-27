import crypto from 'crypto'

type TConfig = { cipherAlgo: string; digestAlgo: string; encoding: 'utf8' }

export default class Cipher {
  constructor(
    private readonly key: string,
    private readonly config: TConfig = { cipherAlgo: 'aes192', digestAlgo: 'sha256', encoding: 'utf8' }
  ) {}

  encrypt = (message: string) => {
    const cipher = crypto.createCipher(this.config.cipherAlgo, this.key)
    return [cipher.update(message, this.config.encoding, 'hex'), cipher.final('hex')].join('')
  }

  decrypt = (encryptedMessage: string) => {
    const decipher = crypto.createDecipher(this.config.cipherAlgo, this.key)
    return decipher.update(encryptedMessage, 'hex', this.config.encoding) + decipher.final(this.config.encoding)
  }

  digest = (message: string) => {
    return crypto
      .createHmac(this.config.digestAlgo, this.key)
      .update(new Buffer(message, this.config.encoding))
      .digest('hex')
  }
}
