import Bearer from '../../lib/bearer'

const clientId = 'a-client-id'

describe('function request', () => {
  const instance = new Bearer(clientId)
  //@ts-ignore
  fetch.mockResponse(
    JSON.stringify({ data: [{ uuid: 'patrick', asset: 'patrick-url' }, { uuid: 'something', asset: 'something-url' }] })
  )
  describe('when successful request', () => {
    it.skip('returns payload if success', async () => {
      await instance.functionFetch('integration', 'function')
    })

    it.skip('forwards params', async () => {
      await instance.functionFetch('integration', 'function', {
        query: { something: 'query' },
        somethingElse: 'query'
      })
    })

    describe('with returned error', () => {
      it.skip('returns error payload', async () => {
        await instance.functionFetch('integration', 'function-with-error')
      })
    })
  })
  describe('with server error', () => {
    it.skip('throw an error', async () => {
      await instance.functionFetch('integration', 'server-error')
    })
  })
})
