import { runTransformersOn } from '../utils/helpers'
import transfomer from '../../src/transformers/component-tag-name-scoping'
import Metadata from '../../src/metadata'
import gatherMetadata from '../../src/transformers/gather-metadata'
const TEST_NAME = 'bearer-component-tag-scoping'

describe('without scope', () => {
  const metadata = new Metadata()
  runTransformersOn('scopeless', TEST_NAME, [gatherMetadata({ metadata }), transfomer({ metadata })])
})

describe('With scope', () => {
  const metadata = new Metadata('bearer', 'xyz')
  runTransformersOn('scoped components', TEST_NAME, [gatherMetadata({ metadata }), transfomer({ metadata })])
})
