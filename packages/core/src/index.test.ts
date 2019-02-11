import Bearer, { scopedT, scopedP } from './index'

describe('exports', () => {
  it('has a default', () => {
    expect(Bearer).toBeInstanceOf(Function)
  })

  it(' exports t i18n helper', () => {
    expect(scopedT).toBeInstanceOf(Function)
  })

  it(' exports p i18n helper', () => {
    expect(scopedP).toBeInstanceOf(Function)
  })
})
