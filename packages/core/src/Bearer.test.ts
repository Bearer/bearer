import Bearer from './Bearer'

describe('Bearer', () => {
  describe('Init', () => {
    it('Returns singleton (exact same instance)', () => {
      expect(Bearer.config).toBe(Bearer.config)
    })

    it('Set default values ', () => {
      expect(Bearer.config.integrationHost).toEqual(
        'https://int.bearer.sh/'
      )
      expect(Bearer.config.loadingComponent).toBeUndefined()
    })

    it('Overrides defaults ', () => {
      Bearer.init({ integrationHost: 'spongebob' })
      expect(Bearer.config.integrationHost).toEqual('spongebob')
    })
  })
})
