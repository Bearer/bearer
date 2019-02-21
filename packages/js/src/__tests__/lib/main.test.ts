import bearer from '../../lib/main'

describe('main', () => {
  it('exports a function expecting 1 argument', () => {
    expect(bearer).toBeInstanceOf(Function)
    expect(bearer).toHaveLength(1)
  })
})
