import Cipher from '../src/index'

describe('Cipher', () => {
  it('encrypts and decrypts the data correctly', async () => {
    expect.assertions(1)
    const message = 'test message'
    const cipher = new Cipher({ key: 'test' })
    const processed = await cipher.decrypt(cipher.encrypt(message))
    return expect(processed).toBe(message)
  })
})
