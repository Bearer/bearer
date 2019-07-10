import nock from 'nock'

import clientFactory, { Bearer } from './client'
const apiKey = 'spongeBobApiKey'
const okResponse = { ok: 'ok' }
const distantApi = jest.fn(() => okResponse)

describe('Bearer client', () => {
  const client = clientFactory(apiKey)

  it('returns a client instance', () => {
    expect(client).toBeInstanceOf(Bearer)
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
      expect(data).toEqual(okResponse)
    })
  })

  describe('#integration', () => {
    const integrationName = '12345'
    const api = client.integration(integrationName)

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
      expect(data).toEqual(okResponse)
    })

    it('allows to make authenticated API calls', async () => {
      const authId = 'abcde12345...'
      distantApi.mockClear()
      nock('https://int.bearer.sh', {
        reqheaders: {
          authorization: apiKey,
          'Bearer-Auth-Id': authId
        }
      })
        .post(`/api/v4/functions/backend/${integrationName}/bearer-proxy/test`)
        .query({ sponge: 'bob' })
        .reply(200, distantApi)

      const { data } = await api.auth(authId).post('/test', { query: { sponge: 'bob' } })

      expect(distantApi).toHaveBeenCalled()
      expect(data).toEqual(okResponse)
    })

    it('has an alias function "authenticate"', async () => {
      expect(api.authenticate).toEqual(api.auth)
    })

    it('sends any configured setup id in a Bearer-Setup-Id header', async () => {
      const setupId = 'test-setup-id'

      nock('https://int.bearer.sh', {
        reqheaders: {
          authorization: apiKey,
          'Bearer-Setup-Id': setupId
        }
      })
        .get(`/api/v4/functions/backend/${integrationName}/bearer-proxy/test`)
        .reply(200, distantApi)

      const { data } = await api.setup(setupId).get('/test')

      expect(distantApi).toHaveBeenCalled()
      expect(data).toEqual(okResponse)
    })
  })
})
