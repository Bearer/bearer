import crypto from 'crypto'

type TConfig = { cipherAlgo: string; digestAlgo: string; encoding: 'utf8' }
const IV_LENGTH = 16 // For AES, this is always 16

export default class Cipher {
  constructor(
    private readonly key: string,
    private readonly config: TConfig = { cipherAlgo: 'aes-256-cbc', digestAlgo: 'sha256', encoding: 'utf8' }
  ) {
    if (!key || !key.length) {
      throw new Error('Invalid key length')
    }
  }

  encrypt = (message: string) => {
    const iv = crypto.randomBytes(IV_LENGTH)
    const cipher = crypto.createCipheriv(this.config.cipherAlgo, this.key, iv)
    const encrypted = Buffer.concat([cipher.update(message), cipher.final()])

    return [iv.toString('hex'), encrypted.toString('hex')].join(':')
  }

  decrypt = (encryptedMessage: string) => {
    const textParts = encryptedMessage.split(':')
    const iv = Buffer.from(textParts.shift()!, 'hex')
    const encryptedText = Buffer.from(textParts.join(':'), 'hex')
    const decipher = crypto.createDecipheriv('aes-256-cbc', Buffer.from(this.key), iv)
    const decrypted = decipher.update(encryptedText)

    return Buffer.concat([decrypted, decipher.final()]).toString()
  }

  digest = (message: string) => {
    return crypto
      .createHmac(this.config.digestAlgo, this.key)
      .update(Buffer.from(message, this.config.encoding))
      .digest('hex')
  }
}
