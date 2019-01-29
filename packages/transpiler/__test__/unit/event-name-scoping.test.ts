import transfomer, { isBearerEvent } from '../../src/transformers/event-name-scoping'
import { runTransformersOn } from '../utils/helpers'
import Metadata from '../../src/metadata'

const TEST_NAME = 'event-name-scoping'

describe('Scope Event Names', () => {
  const metadata = new Metadata()
  runTransformersOn('transformer', TEST_NAME, [transfomer({ metadata })])
})

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
