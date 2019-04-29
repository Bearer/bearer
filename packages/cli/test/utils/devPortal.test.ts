import { devPortalClient } from '../../src/utils/devPortal'
import * as nock from 'nock'

describe('devPortalClient', () => {
  it('export a facotry and expect a command as argument', () => {
    expect(devPortalClient).toHaveLength(1)
  })

  it('creates an instance with a request async method', () => {
    const instance = devPortalClient({ constants: { DeveloperPortalAPIUrl: 'http://test.bearer.sh' } } as any)
    expect(instance.request).toHaveLength(1)
  })

  describe('instance', () => {
    function setup({ token }: { token?: { access_token: string } }) {
      const command = {
        constants: { DeveloperPortalAPIUrl: 'http://test.bearer.sh' },
        bearerConfig: { getToken: jest.fn(() => token) }
      } as any
      const instance = devPortalClient(command)
      return instance
    }

    it('throws error no token present', async () => {
      nock('http://test.bearer.sh', { reqheaders: { authorization: 'Bearer valid_token' } })
        .post('/', { query: 'somehting' })
        .reply(200, 'OK')
      const instance = setup({ token: { access_token: 'valid_token' } })

      const { data } = await instance.request({ query: 'somehting' })

      expect(data).toEqual('OK')
    })
  })
})
