import transfomer from '../../src/transformers/i18n-modifier'
import { runTransformersOn } from '../utils/helpers'

const TEST_NAME = 'i18n-modifier'

describe('Scope Event Names', () => {
  runTransformersOn('transformer', TEST_NAME, [transfomer({})])
})
