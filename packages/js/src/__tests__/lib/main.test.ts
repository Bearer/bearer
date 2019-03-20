import bearer from '../../lib/main'
import { I18n } from '../../lib/i18n'

describe('main', () => {
  it('exports a function expecting 2 arguments', () => {
    expect(bearer).toBeInstanceOf(Function)
    expect(bearer).toHaveLength(2)
  })

  it('forwards options to instance', () => {
    const options = {
      domObserver: false,
      integrationHost: 'something',
      refreshDebounceDelay: 0
    }
    const instance = bearer('token', options)

    expect(instance.config).toMatchObject(expect.objectContaining(options))
  })

  it('use defaults', () => {
    const instance = bearer('token', { integrationHost: undefined })

    expect(instance.config).toMatchObject(
      expect.objectContaining({
        domObserver: true,
        integrationHost: 'INTEGRATION_HOST_URL',
        refreshDebounceDelay: 200
      })
    )
  })

  it('has a i18n accessor', () => {
    expect(bearer.i18n).toBeInstanceOf(I18n)
  })

  it('has a secured accessor', () => {
    expect(bearer.secured).toBeFalsy()
  })

  it('has a secured setter', () => {
    bearer.secured = true
    expect(bearer.secured).toBeTruthy()
  })

  describe('instance', () => {
    beforeEach(() => {
      // @ts-ignore
      bearer.instance = undefined
    })

    it('generates a new instance with a null client id', () => {
      expect(bearer.instance.clientId).toBe(undefined)
    })

    it('use the one setup previously', () => {
      bearer('existing-client-id')
      expect(bearer.instance.clientId).toBe('existing-client-id')
    })
  })
})
