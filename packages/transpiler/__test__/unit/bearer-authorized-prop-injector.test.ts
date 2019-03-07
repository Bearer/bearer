import { runTransformersOn } from '../utils/helpers'
import transfomer from '../../src/transformers/bearer-authorized-integration-id-prop-injector'

const TEST_NAME = 'bearer-authorized-prop-injector'

describe('bearer-authorized required property', () => {
  runTransformersOn('transformer', TEST_NAME, [transfomer({})])
})
