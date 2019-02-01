import i18n from './index'
import { I18nStore } from './store'

describe('i18n', () => {
  it('export a function expecting a store', () => {
    expect(i18n).toBeInstanceOf(Function)
    expect(i18n).toHaveLength(1)
  })

  it('returns a function', () => {
    const result = i18n(null)
    expect(result).toBeInstanceOf(Function)
    expect(result).toHaveLength(3)
  })

  describe('custom store', () => {
    const dictionnary = {
      interpolatesomething: '{{yeah}} !!!',
      'my.key': 'existing value'
    }
    const store: I18nStore = {
      get: (key: string) => dictionnary[key],
      setLocale: jest.fn(),
      loadLocale: jest.fn()
    }
    const instance = i18n(store)

    it('returns existing key', () => {
      expect(instance('my.key', 'Default Value')).toEqual('existing value')
    })

    it('returns existing key', () => {
      expect(instance('interpolatesomething', 'Default Value', { yeah: 'wonderful' })).toEqual('wonderful !!!')
    })

    it('returns default value', () => {
      expect(instance('missing.key', 'Default Value')).toEqual('Default Value')
    })
  })
})
