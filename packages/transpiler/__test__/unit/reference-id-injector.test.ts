import { runTransformersOn } from '../utils/helpers'
import transformer from '../../src/transformers/reference-id-injector'

const TEST_NAME = 'reference-id-injector'

describe('reference id injector', () => {
  runTransformersOn('change files content', TEST_NAME, [transformer()])
})
