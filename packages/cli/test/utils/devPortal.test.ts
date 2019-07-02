import { devPortalClient, withFreshToken } from '../../src/utils/devPortal'

jest.mock('../../src/actions/login')

import { promptToLogin } from '../../src/actions/login'

describe('devPortalClient', () => {
  it('export a factory and expect a command as argument', () => {
    expect(devPortalClient).toHaveLength(1)
  })

  it('creates an instance with a request async method', () => {
    const instance = devPortalClient({ constants: { DeveloperPortalAPIUrl: 'http://test.bearer.sh' } } as any)
    expect(instance.request).toHaveLength(1)
  })
})

describe('withFreshToken', () => {
  describe('when token is still valid', () => {
    it('does nothing', async () => {
      const command = {
        bearerConfig: {
          getToken: jest.fn(() => ({
            expires_at: Date.now() + 2000,
            refresh_token: 'a token'
          }))
        }
      }
      await withFreshToken(command as any)
    })
  })

  describe('when token does not exists (never logged in)', () => {
    it('prompt to log in', async () => {
      const command = {
        bearerConfig: {
          getToken: jest.fn()
        },
        colors: {
          bold: jest.fn()
        },
        log: jest.fn()
      }

      await withFreshToken(command as any)

      expect(promptToLogin).toHaveBeenCalledTimes(1)
    })
  })

  describe('when token is outdated', () => {
    it('refresh the token', () => {})
  })
})
