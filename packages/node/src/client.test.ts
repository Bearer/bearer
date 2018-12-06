import nock from 'nock'

import clientFactory, { BearerClient } from './client'
const clientId = 'spongeBobClientId'

describe('Bearer client', () => {
  const client = clientFactory(clientId)

  it('returns a client instance', () => {
    expect(client).toBeInstanceOf(BearerClient)
  })

  describe('#call', () => {
    it('send request to the intent', async () => {
      const distantApi = jest.fn(() => ({ ok: 'ok' }))

      nock('https://int.bearer.sh', {
        reqheaders: {
          authorization: clientId
        }
      })
        .post('/backend/api/v1/12345-scenario-name/intentName')
        .reply(200, distantApi)

      const { data } = await client.call('12345-scenario-name', 'intentName')

      expect(distantApi).toHaveBeenCalled()
      expect(data).toEqual({ ok: 'ok' })
    })
  })
})
