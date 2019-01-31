import defaultStore, { Store } from './store'

describe('i18n store', () => {
  it('exports a default instance', () => {
    expect(defaultStore).toBeInstanceOf(Store)
  })
})
