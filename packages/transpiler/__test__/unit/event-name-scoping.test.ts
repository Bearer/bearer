import { isBearerEvent } from '../../src/transformers/event-name-scoping'

// import { runUnitOn } from '../utils/helpers'
// const TEST_NAME = 'event-name-scoping'

// describe('Scope Event Names', () => {
//   runUnitOn(TEST_NAME)
// })

describe('isBearerEvent', () => {
  describe('when non matching event name', () => {
    it('does not match', () => {
      expect(isBearerEvent('anEvent')).toBeFalsy()
    })

    it('does not match body scoped events', () => {
      expect(isBearerEvent('body:anEvent')).toBeFalsy()
    })
  })
  describe('when matching event name', () => {
    it('matches', () => {
      expect(isBearerEvent('bearer-something')).toBeTruthy()
    })

    it('matches body scoped events', () => {
      expect(isBearerEvent('body:bearer-something')).toBeTruthy()
    })
  })
})
