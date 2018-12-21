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

  it('returns a digest', async () => {
    const expectedDigest = '27764562547d0665075c1fcf972ee5990db168ee87e6888b4aed2ba2c0f3085d'
    const body = {
      key1: 'value1',
      key2: 'value2'
    }
    const digest = await cipher.digest(JSON.stringify(body))
    expect(digest).toBe(expectedDigest)
  })
})
