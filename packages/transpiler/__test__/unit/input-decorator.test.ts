import { runTransformersOn } from '../utils/helpers'
import Metadata from '../../src/metadata'
import GatherMetadata from '../../src/transformers/gather-metadata'
import InputDecorator from '../../src/transformers/input-decorator'
const TEST_NAME = 'input-decorator'

const metadata = new Metadata('bearer', 'test')

describe('Input decorator units', () => {
  runTransformersOn('generates stuff', TEST_NAME, [GatherMetadata({ metadata }), InputDecorator({ metadata })])
})
