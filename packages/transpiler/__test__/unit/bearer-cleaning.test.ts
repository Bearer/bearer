import { runTransformersOn } from '../utils/helpers'
import transformer from '../../src/transformers/bearer-cleaning'
const TEST_NAME = 'bearer-cleaning'

describe('bearer cleaning transformer', () => {
  runTransformersOn('transformer', TEST_NAME, [transformer()])
})
