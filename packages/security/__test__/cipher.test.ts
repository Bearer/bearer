import Cipher from '../src/index'
const cipher = new Cipher('12345678901234567890123456789012')

describe('Cipher', () => {
  it('raises an error if key length is incorrect', () => {
    expect(() => {
      new Cipher('')
    }).toThrowError('Invalid key length')
  })

  it('encrypts the data correctly', () => {
    expect.assertions(2)
    const message = 'test message'

    const encrypted = cipher.encrypt(message)

    expect(encrypted).not.toBe(message)
    expect(encrypted).not.toEqual(cipher.encrypt(message))
  })

  it('decrypts the data correctly', () => {
    const encrypted = 'b77ab4b52084ee0fb7dbfd30f1379894:b9bd456f63dcb943d0600670d346910a'

    const decrypted = cipher.decrypt(encrypted)

    expect(decrypted).toBe('test message')
  })

  it('encrypt and decrypt correctly', () => {
    const message = 'test message'

    expect(cipher.decrypt(cipher.encrypt(message))).toBe('test message')
  })

  it('returns a digest', () => {
    const expectedDigest = 'd0f7c48b77ef540f8704f90eec06d836da0b72ff7db779755dbee1b08c9781b9'

    const body = {
      key1: 'value1',
      key2: 'value2'
    }
    const digest = cipher.digest(JSON.stringify(body))

    expect(digest).toBe(expectedDigest)
  })
})
