import Bearer, { t } from './index'

describe('exports', () => {
  it('has a default', () => {
    expect(Bearer).toBeInstanceOf(Function)
  })

  it(' exports i18n helper', () => {
    expect(t).toBeInstanceOf(Function)
  })
})
