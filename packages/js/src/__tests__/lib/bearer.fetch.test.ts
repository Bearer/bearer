import Bearer from '../../lib/bearer'

const clientId = 'a-client-id'

describe('function request', () => {
  const instance = new Bearer(clientId)

  describe('when successful request', () => {
    it('returns payload if success', async () => {
      await instance.functionFetch('integration', 'function')
    })

    it('forwards params', async () => {
      await instance.functionFetch('integration', 'function', {
        query: { something: 'query' },
        somethingElse: 'query'
      })
    })

    describe('with returned error', () => {
      it('returns error payload', async () => {
        await instance.functionFetch('integration', 'function-with-error')
      })
    })
  })
  describe('with server error', () => {
    it('throw an error', async () => {
      await instance.functionFetch('integration', 'server-error')
    })
  })
})
