import { runTransformersOn } from '../utils/helpers'

import Metadata from '../../src/metadata'
import GatherMetadata from '../../src/transformers/gather-metadata'
import OutputDecorator from '../../src/transformers/output-decorator'
const TEST_NAME = 'input-output-intent-arguments'

const metadata = new Metadata('bearer', 'test')

describe('input and output arguments', () => {
  runTransformersOn('generates the refIds properly', TEST_NAME, [
    GatherMetadata({ metadata }),
    OutputDecorator({ metadata })
  ])
})
