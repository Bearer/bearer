import { runTransformersOn } from '../utils/helpers'
import Metadata from '../../src/metadata'
import GatherMetadata from '../../src/transformers/gather-metadata'
import transformer from '../../src/transformers/output-decorator'

const metadata = new Metadata('bearer', 'test')
const TEST_NAME = 'output-decorator'

describe('Output decorator units', () => {
  runTransformersOn('generates stuff', TEST_NAME, [GatherMetadata({ metadata }), transformer({ metadata })])
})
