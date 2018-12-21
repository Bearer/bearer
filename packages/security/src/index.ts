import crypto from 'crypto'

export default class Cipher {
  constructor(private readonly config: { key: string }) {}

  encrypt = (message: string) => {
    const cipher = crypto.createCipher('aes192', this.config.key)

    return [cipher.update(message, 'utf8', 'hex'), cipher.final('hex')].join('')
  }

  decrypt = (encryptedMessage: string) => {
    return new Promise((resolve, reject) => {
      let decrypted = ''
      const decipher = crypto.createDecipher('aes192', this.config.key)
      decipher.on('readable', () => {
        const data = decipher.read()
        if (data) {
          decrypted += data.toString()
        }
      })
      decipher.on('end', () => {
        resolve(decrypted)
      })

      decipher.on('error', e => {
        reject(e)
      })

      try {
        decipher.write(encryptedMessage, 'hex')
      } catch (e) {
        console.log(`Malformed signature: ${e.message}`)
        reject(e)
      }
      decipher.end()
    })
  }

  digest = (message: string) => {
    return crypto
      .createHmac('sha256', this.config.key)
      .update(new Buffer(message, 'utf-8'))
      .digest('hex')
  }
}
