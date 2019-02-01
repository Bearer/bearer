import Bearer, { t, p } from './index'

describe('exports', () => {
  it('has a default', () => {
    expect(Bearer).toBeInstanceOf(Function)
  })

  it(' exports t i18n helper', () => {
    expect(t).toBeInstanceOf(Function)
  })

  it(' exports p i18n helper', () => {
    expect(p).toBeInstanceOf(Function)
  })
})
