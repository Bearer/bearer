import nock from 'nock'

import clientFactory, { BearerClient, ScenarioClient } from './client'
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
        .post('/backend/api/v1/12345-scenario-name/intentName')
        .reply(200, distantApi)

      const { data } = await client.call('12345-scenario-name', 'intentName')

      expect(distantApi).toHaveBeenCalled()
      expect(data).toEqual({ ok: 'ok' })
    })
  })
})

describe('ScenarioClient', () => {
  const token = 'a-different-token'
  const anotherScenarioName = 'scenario-name'
  const client = new ScenarioClient(token, {}, anotherScenarioName)

  it('creates a scenario client', () => {
    expect(client).toBeInstanceOf(ScenarioClient)
  })

  it('calls correct scenario intents', async () => {
    distantApi.mockClear()
    nock('https://int.bearer.sh', {
      reqheaders: {
        authorization: token
      }
    })
      .post(`/backend/api/v1/${anotherScenarioName}/intentName`)
      .query({ sponge: 'bob' })
      .reply(200, distantApi)

    const { data } = await client.call('intentName', { query: { sponge: 'bob' } })

    expect(distantApi).toHaveBeenCalled()
    expect(data).toEqual({ ok: 'ok' })
  })
})
