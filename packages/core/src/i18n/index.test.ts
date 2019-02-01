import { translate, pluralize } from './index'
import { I18nStore, Store } from './store'

describe('i18n.translate', () => {
  it('export a function expecting a store', () => {
    expect(translate).toBeInstanceOf(Function)
    expect(translate).toHaveLength(1)
  })

  it('returns a function', () => {
    const result = translate(null)
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
    const instance = translate(store)

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

describe('i18n.pluralize', () => {
  it('export a function expecting a store', () => {
    expect(pluralize).toBeInstanceOf(Function)
    expect(pluralize).toHaveLength(1)
  })

  it('returns a function', () => {
    const result = pluralize(null)
    expect(result).toBeInstanceOf(Function)
    expect(result).toHaveLength(4)
  })

  describe('custom store', () => {
    const dictionnary = {
      en: {
        my: {
          key: {
            0: 'none',
            1: 'one {{yeah}}',
            2: 'two',
            many: 'the rest'
          }
        }
      }
    }
    const store = new Store(dictionnary)
    const instance = pluralize(store)

    it('returns existing key', () => {
      expect(instance('my.key', 0, 'Default Value')).toEqual('none')
    })

    it('returns existing key interpolated', () => {
      expect(instance('my.key', 1, 'Default Value', { yeah: 'wonderful' })).toEqual('one wonderful')
    })

    it('defined quantity', () => {
      expect(instance('my.key', 2, 'Default Value')).toEqual('two')
    })

    it('use many', () => {
      expect(instance('my.key', 3, 'Default Value')).toEqual('the rest')
    })

    it('return default', () => {
      expect(instance('unknoown', 2, 'Default Value')).toEqual('Default Value')
    })
  })
})
