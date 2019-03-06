import nock from 'nock'

import clientFactory, { BearerClient, IntegrationClient } from './client'
const clientId = 'spongeBobClientId'

const distantApi = jest.fn(() => ({ ok: 'ok' }))

describe('Bearer client', () => {
  const client = clientFactory(clientId)

  it('returns a client instance', () => {
    expect(client).toBeInstanceOf(BearerClient)
  })

  describe('#call', () => {
    it('send request to the intent', async () => {
      distantApi.mockClear()
      nock('https://int.bearer.sh', {
        reqheaders: {
          authorization: clientId
        }
      })
        .post('/api/v3/intents/backend/12345-integration-name/intentName')
        .reply(200, distantApi)

      const { data } = await client.call('12345-integration-name', 'intentName')

      expect(distantApi).toHaveBeenCalled()
      expect(data).toEqual({ ok: 'ok' })
    })
  })
})

describe('IntegrationClient', () => {
  const token = 'a-different-token'
  const anotherIntegrationName = 'integration-name'
  type TIntegrationIntentNames = 'intent-name' | 'other-intent'
  const client = new IntegrationClient<TIntegrationIntentNames>(token, {}, anotherIntegrationName)

  it('creates a integration client', () => {
    expect(client).toBeInstanceOf(IntegrationClient)
  })

  it('calls correct integration intents', async () => {
    distantApi.mockClear()
    nock('https://int.bearer.sh', {
      reqheaders: {
        authorization: token
      }
    })
      .post(`/api/v3/intents/backend/${anotherIntegrationName}/intent-name`)
      .query({ sponge: 'bob' })
      .reply(200, distantApi)

    const { data } = await client.call('intent-name', { query: { sponge: 'bob' } })

    expect(distantApi).toHaveBeenCalled()
    expect(data).toEqual({ ok: 'ok' })
  })
})
