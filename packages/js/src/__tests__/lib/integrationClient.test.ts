import nock from 'nock'

import { IntegrationClient } from '../../lib/integrationClient'

describe('IntegrationClient', () => {
  const integrationId = 'integration-id'
  const instance = new IntegrationClient({ config: { integrationHost: 'https://int.bearer.sh' } } as any, integrationId)

  describe('constructor', () => {
    it('set properties', () => {
      // @ts-ignore
      expect(instance.bearerInstance).toBeDefined()
      expect(instance.integrationId).toBeDefined()
      expect(instance.authId).not.toBeDefined()
      expect(instance.setupId).not.toBeDefined()
    })
  })

  describe('#auth', () => {
    it('creates a new instance with the given authId', () => {
      const withAuth = instance.auth('auth-id')

      expect(withAuth).not.toEqual(instance)
      expect(withAuth.authId).toEqual('auth-id')
    })
  })

  describe('#setup', () => {
    it('creates a new instance with the given setupId', () => {
      const withSetup = instance.setup('setup-id')

      expect(withSetup).not.toEqual(instance)
      expect(withSetup.setupId).toEqual('setup-id')
    })
  })

  describe('#invoke', () => {
    const apiResponse = jest.fn(() => ({ data: 'ok' }))

    it('performs a post request', async () => {
      nock('https://int.bearer.sh', {})
        .intercept(`/api/v4/functions/${integrationId}/customFunction`, 'post', { bodyData: 'ok' })
        .once()
        .reply(200, apiResponse)

      const response = await instance.invoke('customFunction', { bodyData: 'ok' })

      expect(apiResponse).toHaveBeenCalled()
      expect(response).toEqual({ data: 'ok' })
    })

    it('forwards setup and auth', async () => {
      nock('https://int.bearer.sh', {})
        .intercept(`/api/v4/functions/${integrationId}/customFunction`, 'post', { extraData: { key: 'value' } })
        .once()
        .query({
          authId: 'my-auth-id',
          setupId: 'my-setup-id'
        })
        .reply(200, apiResponse)

      const instanceWithOptions = instance.auth('my-auth-id').setup('my-setup-id')
      const response = await instanceWithOptions.invoke('customFunction', { extraData: { key: 'value' } })

      expect(apiResponse).toHaveBeenCalled()
      expect(response).toEqual({ data: 'ok' })
    })
  })

  describe('proxy', () => {
    const okResponse = { ok: 'ok' }
    const distantApi = jest.fn(() => okResponse)
    const headers = { extra: 'headers' }
    const query = { extraQuery: 'value' }

    function mockRequest({ method, body }: any) {
      nock('https://int.bearer.sh', {})
        .intercept(`/api/v4/functions/${integrationId}/bearer-proxy/test`, method, body)
        .once()
        .query(query)
        .reply(200, distantApi)
    }

    beforeEach(() => {
      distantApi.mockClear()
    })

    describe('#get', () => {
      it('performs GET request', async () => {
        mockRequest({ method: 'GET' })

        const { data } = await instance.get('/test', { headers, query })

        expect(distantApi).toHaveBeenCalled()
        expect(data).toEqual(okResponse)
      })
    })

    describe('#head', () => {
      it('performs HEAD request', async () => {
        mockRequest({ method: 'HEAD' })

        const { data } = await instance.head('/test', { headers, query })

        expect(distantApi).toHaveBeenCalled()
        expect(data).toEqual(okResponse)
      })
    })

    describe('#post', () => {
      it('performs POST request', async () => {
        mockRequest({ body: 'POST body', method: 'POST' })

        const { data } = await instance.post('/test', { headers, query, body: 'POST body' })

        expect(distantApi).toHaveBeenCalled()
        expect(data).toEqual(okResponse)
      })
    })

    describe('#put', () => {
      it('performs PUT request', async () => {
        mockRequest({ body: { json: 'PUT body' }, method: 'PUT' })

        const { data } = await instance.put('/test', { headers, query, body: { json: 'PUT body' } })

        expect(distantApi).toHaveBeenCalled()
        expect(data).toEqual(okResponse)
      })
    })

    describe('#patch', () => {
      it('performs PATCH request', async () => {
        mockRequest({ body: 'PATCH body', method: 'PATCH' })

        const { data } = await instance.patch('/test', { headers, query, body: 'PATCH body' })

        expect(distantApi).toHaveBeenCalled()
        expect(data).toEqual(okResponse)
      })
    })

    describe('#delete', () => {
      it('performs DELETE request', async () => {
        mockRequest({ body: 'DELETE body', method: 'DELETE' })

        const { data } = await instance.delete('/test', { headers, query, body: 'DELETE body' })

        expect(distantApi).toHaveBeenCalled()
        expect(data).toEqual(okResponse)
      })
    })
  })
})
