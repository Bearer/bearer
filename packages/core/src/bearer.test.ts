import Bearer from './bearer'

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

    describe('#askAuthorizations', () => {
      it('opens popup with  the correct url', () => {
        // @ts-ignore
        global.bearer = { clientId: 'askAuthorizations-clientId' }
        const win = { open: jest.fn() }
        const instance = new Bearer({ integrationHost: 'https://trash.bearer.sh/', secured: true }, win as any)
        // @ts-ignore
        instance.sessionInitialized()
        expect(instance.askAuthorizations({ scenarioId: 'ok', setupId: 'ok', authRefId: 'IAM' })).toBeTruthy()
        expect(win.open).toHaveBeenCalledWith(
          'https://trash.bearer.sh/v2/auth/ok?setupId=ok&authId=IAM&secured=true&clientId=askAuthorizations-clientId',
          '',
          'resizable,scrollbars,status,centerscreen=yes,width=500,height=600'
        )
      })

      it('does not open ', () => {
        const instance = Bearer.init()
        expect(instance.askAuthorizations({ scenarioId: 'ok', setupId: 'ok' })).toBeFalsy()
      })
    })
  })
})
