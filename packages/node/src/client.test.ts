import nock from 'nock'

import clientFactory, { BearerClient, IntegrationClient } from './client'
const clientId = 'spongeBobClientId'

const distantApi = jest.fn(() => ({ ok: 'ok' }))

describe('Bearer client', () => {
  const client = clientFactory(clientId)

  it('returns a client instance', () => {
    expect(client).toBeInstanceOf(BearerClient)
  })

  describe('#invoke', () => {
    it('send request to the function', async () => {
      distantApi.mockClear()
      nock('https://int.bearer.sh', {
        reqheaders: {
          authorization: clientId
        }
      })
        .post('/api/v4/functions/backend/12345-integration-name/functionName')
        .reply(200, distantApi)

      const { data } = await client.invoke('12345-integration-name', 'functionName')

      expect(distantApi).toHaveBeenCalled()
      expect(data).toEqual({ ok: 'ok' })
    })
  })
})

describe('IntegrationClient', () => {
  const token = 'a-different-token'
  const anotherIntegrationName = 'integration-name'
  type TIntegrationFunctionNames = 'function-name' | 'other-function'
  const client = new IntegrationClient<TIntegrationFunctionNames>(token, {}, anotherIntegrationName)

  it('creates a integration client', () => {
    expect(client).toBeInstanceOf(IntegrationClient)
  })

  it('invokes correct integration functions', async () => {
    distantApi.mockClear()
    nock('https://int.bearer.sh', {
      reqheaders: {
        authorization: token
      }
    })
      .post(`/api/v4/functions/backend/${anotherIntegrationName}/function-name`)
      .query({ sponge: 'bob' })
      .reply(200, distantApi)

    const { data } = await client.invoke('function-name', { query: { sponge: 'bob' } })

    expect(distantApi).toHaveBeenCalled()
    expect(data).toEqual({ ok: 'ok' })
  })
})
