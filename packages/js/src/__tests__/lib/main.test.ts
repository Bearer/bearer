import bearer from '../../lib/main'
import { I18n } from '../../lib/i18n'

describe('main', () => {
  it('exports a function expecting 1 argument', () => {
    expect(bearer).toBeInstanceOf(Function)
    expect(bearer).toHaveLength(1)
  })

  it('has a i18n accessor', () => {
    expect(bearer.i18n).toBeInstanceOf(I18n)
  })
})
