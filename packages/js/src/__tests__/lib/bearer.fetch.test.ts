import Bearer from '../../lib/bearer'

const clientId = 'a-client-id'
const returnedData = [{ uuid: 'patrick', asset: 'patrick-url' }, { uuid: 'something', asset: 'something-url' }]

describe('function request', () => {
  const instance = new Bearer(clientId, { integrationHost: 'custom-host' })

  beforeEach(() => {
    // @ts-ignore
    fetch.resetMocks()
  })
  describe('when successful request', () => {
    it('returns payload if success', async () => {
      //@ts-ignore
      fetch.once(JSON.stringify({ data: returnedData }))

      const response = await instance.functionFetch('integration', 'function')

      expect(response.data).toEqual(returnedData)
      expect(response.referenceId).toBe(null)
    })

    it('sets referenceId from meta', async () => {
      // @ts-ignore
      fetch.once(JSON.stringify({ data: returnedData, meta: { referenceId: 'a-reference' } }))

      const response = await instance.functionFetch('integration', 'function', {
        query: { something: 'query' },
        somethingElse: 'query'
      })

      expect(response.data).toEqual(returnedData)
      expect(response.referenceId).toBe('a-reference')
    })

    it('forwards params', async () => {
      // @ts-ignore
      fetch.once(JSON.stringify({}))

      await instance.functionFetch('integration', 'function', {
        query: { something: 'query' },
        somethingElse: 'query'
      })

      expect(fetch).toHaveBeenCalledWith(
        'custom-host/api/v3/functions/integration/function?something=query&clientId=a-client-id',
        {
          body: '{"somethingElse":"query"}',
          credentials: 'include',
          headers: { 'content-type': 'application/json' },
          method: 'POST'
        }
      )
    })

    describe('with returned error', () => {
      it('returns error payload', async () => {
        // @ts-ignore
        fetch.once(JSON.stringify({ error: 'something' }))

        expect(instance.functionFetch('integration', 'function-with-error')).rejects.toEqual({
          // 2 levels of error because of jest
          error: { error: 'something' }
        })
      })
    })
  })
  describe('with server error', () => {
    it('throws an error', async () => {
      //@ts-ignore
      fetch.mockReject(new Error('fake error message'))

      expect(instance.functionFetch('integration', 'server-error')).rejects.toEqual({
        error: new Error('fake error message')
      })
    })
  })
})
