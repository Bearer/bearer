import { randomNames } from '../../src/utils/random'

describe('randomNames', () => {
  it('generates 30 names by default', () => {
    const names = randomNames()
    expect(names).toHaveLength(30)
    expect(typeof names[0]).toBe('string')
  })

  it('generates n names ', () => {
    const names = randomNames(10)
    expect(names).toHaveLength(10)
    expect(typeof names[0]).toBe('string')
  })
})
