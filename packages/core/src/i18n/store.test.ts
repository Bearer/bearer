import defaultStore, { Store } from './store'

describe('i18n store', () => {
  it('exports a default instance', () => {
    expect(defaultStore).toBeInstanceOf(Store)
  })

  describe('complete example', () => {
    it('works like a charm', () => {
      const data = {}
      const instance = new Store(data)
      expect(instance.get('my.key')).toBeUndefined()

      instance.loadLocale('en', { my: { key: 'exists' } })
      expect(instance.get('my.key')).toBe('exists')

      instance.loadLocale('en', { my: { key: 'overriden' } })
      expect(instance.get('my.key')).toBe('overriden')

      instance.loadLocale('fr', { my: { key: 'une valeur' } })
      instance.setLocale('fr')
      expect(instance.get('my.key')).toBe('une valeur')

      expect(data).not.toEqual({})
    })
  })

  describe('.get', () => {
    const data = {
      en: {
        my: {
          key: 'exists'
        }
      },
      fr: {
        my: {
          key: 'existe'
        }
      }
    }

    const instance = new Store(data)

    it('retrieves existing value', () => {
      expect(instance.get('my.key')).toBe('exists')
    })

    it('retrieves existing from the correct locale', () => {
      const instance = new Store(data)
      instance.setLocale('fr')
      expect(instance.get('my.key')).toBe('existe')
    })

    describe('when key path does not exist', () => {
      it('returns undefined', () => {
        expect(instance.get('my.missing.key')).toBeUndefined()
      })
    })
  })
})
