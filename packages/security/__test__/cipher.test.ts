import Cipher from '../src/index'

const cipher = new Cipher({ key: 'test' })
describe('Cipher', () => {
  it('encrypts the data correctly', async () => {
    expect.assertions(1)
    const message = 'test message'
    const encrypted = cipher.encrypt(message)
    return expect(encrypted).toBe('873e48ebd52490bdec6db078b26047f2')
  })

  it('decrypts the data correctly', async () => {
    const encrypted = '873e48ebd52490bdec6db078b26047f2'
    const decrypted = await cipher.decrypt(encrypted)
    return expect(decrypted).toBe('test message')
  })
})
