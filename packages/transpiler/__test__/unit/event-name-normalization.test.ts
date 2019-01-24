import { runTransformersOn } from '../utils/helpers'
import eventNormalize from '../../src/transformers/event-name-normalizer'
const TEST_NAME = 'event-normalization'

describe('normalize event names', () => {
  runTransformersOn('changes event name', TEST_NAME, [eventNormalize()])
})
