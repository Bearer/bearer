import { Authentications } from '@bearer/types/lib/authentications'

import SetupAuth from '../../../src/commands/setup/auth'
import { ensureBearerStructure } from '../../helpers/setup'
import { readFile } from '../../helpers/utils'

jest.mock('open', () => {
  return () => ({
    on: jest.fn()
  })
})

describe('setup', () => {
  describe('auth', () => {
    describe('with argument', () => {
      describe.each([
        [{ auth: Authentications.ApiKey, arg: 'anApikey', returned: {} }],
        [{ auth: Authentications.OAuth1, arg: 'aKey:aConsumerSecret', returned: { someDataFromOauth1Server: 'ok' } }],
        [
          { auth: Authentications.OAuth2, arg: 'aClientId:aClientSecret', returned: { someDataFromOauth2Server: 'ok' } }
        ],
        [{ auth: Authentications.Basic, arg: 'aUsername:aPassword', returned: {} }],
        [{ auth: Authentications.Custom, arg: '', returned: {} }],
        [{ auth: Authentications.NoAuth, arg: '', returned: {} }]
      ])('auth type %j', (...usecase) => {
        const [{ returned, auth, arg }] = usecase
        beforeAll(() => {
          SetupAuth.prototype.fetchAuthToken = () => {
            return Promise.resolve(Buffer.from(JSON.stringify(returned)).toString('base64') as any)
          }
        })

        test(`save to local.config.jsonc`, async () => {
          const path = ensureBearerStructure({
            authConfig: {
              authType: auth
            }
          })

          await SetupAuth.run(['--path', path, arg])
          expect(readSetup(path)).toMatchSnapshot()
        })
      })
    })

    describe('with environment variables', () => {
      const OLD_ENV = process.env

      describe.each([
        [
          {
            auth: Authentications.ApiKey,
            env: { BEARER_AUTH_APIKEY: 'anApiKeyFromEnv' },
            returned: {}
          }
        ],
        [
          {
            auth: Authentications.OAuth1,
            env: { BEARER_AUTH_CONSUMER_KEY: 'aKeyFromEnv', BEARER_AUTH_CONSUMER_SECRET: 'aSecretFromEnv' },
            returned: { someDataFromOauth1Server: 'ok' }
          }
        ],
        [
          {
            auth: Authentications.OAuth2,
            env: { BEARER_AUTH_CLIENT_ID: 'aClientIdFromEnv', BEARER_AUTH_CLIENT_SECRET: 'sSecretFromEmv' },
            returned: { someDataFromOauth2Server: 'ok' }
          }
        ],
        [
          {
            auth: Authentications.Basic,
            env: { BEARER_AUTH_USERNAME: 'userNameFromEnv', BEARER_AUTH_PASSWORD: 'passwordFromEnv' },
            returned: {}
          }
        ],
        [{ auth: Authentications.Custom, env: {}, returned: {} }],
        [{ auth: Authentications.NoAuth, env: {}, returned: {} }]
      ])('%j', (...usecase) => {
        const [{ returned, env, auth }] = usecase
        beforeAll(() => {
          process.env = { ...OLD_ENV, ...env }
          SetupAuth.prototype.fetchAuthToken = () => {
            return Promise.resolve(Buffer.from(JSON.stringify(returned)).toString('base64') as any)
          }

          afterAll(() => {
            process.env = { ...OLD_ENV }
          })
        })

        test(`save to local.config.jsonc`, async () => {
          const path = ensureBearerStructure({
            authConfig: {
              authType: auth
            }
          })

          await SetupAuth.run(['--path', path])
          expect(readSetup(path)).toMatchSnapshot()
        })
      })
    })
  })
})

function readSetup(path: string) {
  return readFile(path, 'local.config.jsonc')
}
