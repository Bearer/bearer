import { i18n } from '@bearer/js'
import { scopedPluralize, scopedTranslate } from './index'

describe('i18n.translate', () => {
  it('export a function expecting a store', () => {
    expect(scopedTranslate).toBeInstanceOf(Function)
    expect(scopedTranslate).toHaveLength(1)
  })

  it('returns a function', () => {
    const result = scopedTranslate(null)
    expect(result).toBeInstanceOf(Function)
    expect(result).toHaveLength(1)
    const translate = result(null)
    expect(translate).toBeInstanceOf(Function)
    expect(translate).toHaveLength(3)
  })

  describe('custom store', () => {
    const store: any = {
      get: jest.fn()
    }

    const instance = scopedTranslate('a-scope')(store)

    it('does a lookup into the store', () => {
      instance('my.key', 'Default Value')
      expect(store.get).toHaveBeenCalledWith('my.key', 'a-scope')
    })
  })
})

describe('i18n.pluralize', () => {
  it('export a function expecting a store', () => {
    expect(scopedPluralize).toBeInstanceOf(Function)
    expect(scopedPluralize).toHaveLength(1)
  })

  it('returns a function', () => {
    const result = scopedPluralize(null)

    expect(result).toBeInstanceOf(Function)
    expect(result).toHaveLength(1)

    const pluralize = result(null)
    expect(pluralize).toBeInstanceOf(Function)
    expect(pluralize).toHaveLength(4)
  })

  describe('custom store', () => {
    const store: any = {
      get: jest.fn()
    }

    const instance = scopedPluralize('a-scope')(store)

    it('does a lookup into the store', () => {
      instance('my.key', 42, 'Default Value', {})
      expect(store.get).toHaveBeenNthCalledWith(1, 'my.key.42', 'a-scope')
      expect(store.get).toHaveBeenLastCalledWith('my.key.many', 'a-scope')
    })
  })
})
