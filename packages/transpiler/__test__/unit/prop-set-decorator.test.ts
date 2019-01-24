import { runTransformersOn } from '../utils/helpers'
import Metadata from '../../src/metadata'
import GatherMetadata from '../../src/transformers/gather-metadata'
import transformer from '../../src/transformers/prop-set-decorator'

const TEST_NAME = 'prop-set-decorator'
const metadata = new Metadata('bearer', 'test')

describe('prop set: hack for react', () => {
  runTransformersOn('apply changes', TEST_NAME, [GatherMetadata({ metadata }), transformer({ metadata })])
})
