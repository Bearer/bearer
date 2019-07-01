import nock from 'nock'

import clientFactory, { BearerClient, BearerClientInstance } from './client'
const apiKey = 'spongeBobApiKey'

const distantApi = jest.fn(() => ({ ok: 'ok' }))

describe('Bearer client', () => {
  const client = clientFactory(apiKey)

  it('returns a client instance', () => {
    expect(client).toBeInstanceOf(BearerClient)
  })

  it('throws an error if the API KEY is not correct', () => {
    expect(() => {
      clientFactory(undefined as any)
    }).toThrowError(
      `Invalid Bearer API key provided.  Value: undefined \
You'll find you API key at this location: https://app.bearer.sh/keys`
    )
  })

  describe('#invoke', () => {
    it('send request to the function', async () => {
      distantApi.mockClear()
      nock('https://int.bearer.sh', {
        reqheaders: {
          authorization: apiKey
        }
      })
        .post('/api/v4/functions/backend/12345-integration-name/functionName')
        .reply(200, distantApi)

      const { data } = await client.invoke('12345-integration-name', 'functionName')

      expect(distantApi).toHaveBeenCalled()
      expect(data).toEqual({ ok: 'ok' })
    })
  })

  describe('#integration', () => {
    const integrationName = '12345'
    const api = client.integration(integrationName)

    it('creates a bearer client instance', () => {
      expect(api).toBeInstanceOf(BearerClientInstance)
    })

    it('performs correct API calls', async () => {
      distantApi.mockClear()
      nock('https://int.bearer.sh', {
        reqheaders: {
          authorization: apiKey
        }
      })
        .post(`/api/v4/functions/backend/${integrationName}/bearer-proxy/test`)
        .query({ sponge: 'bob' })
        .reply(200, distantApi)

      const { data } = await api.post('/test', { query: { sponge: 'bob' } })

      expect(distantApi).toHaveBeenCalled()
      expect(data).toEqual({ ok: 'ok' })
    })
  })
})
