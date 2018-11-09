import Bearer from './Bearer'

describe('Bearer', () => {
  describe('Init', () => {
    it('Returns singleton (exact same instance)', () => {
      expect(Bearer.config).toBe(Bearer.config)
    })

    it('Set default values ', () => {
      expect(Bearer.config.integrationHost).toEqual('BEARER_INTEGRATION_HOST')
      expect(Bearer.config.loadingComponent).toBeUndefined()
    })

    it('Overrides defaults ', () => {
      Bearer.init({ integrationHost: 'spongebob' })
      expect(Bearer.config.integrationHost).toEqual('spongebob')
    })
  })

  describe('Authorization', () => {
    it('with same scenarioId it auhtorizes', done => {
      expect.assertions(1)

      const instance = Bearer.init()
      const callback = jest.fn(() => done())
      Bearer.onAuthorized('scenarioTargeted', callback)

      instance.authorized({ data: { scenarioId: 'scenarioTargeted' } })

      expect(callback).toHaveBeenCalledWith(true)
    })

    it('does not resolve if not a matching scenarioId', done => {
      expect.assertions(2)

      const instance = Bearer.init()

      const callback = jest.fn()
      const otherScenarioCallback = jest.fn(() => done())

      Bearer.onAuthorized('scenarioTargeted', callback)
      Bearer.onAuthorized('otherScenario', otherScenarioCallback)

      instance.authorized({ data: { scenarioId: 'otherScenario' } })

      expect(callback).not.toHaveBeenCalledWith(true)
      expect(otherScenarioCallback).toHaveBeenCalledWith(true)
    })
  })
})
