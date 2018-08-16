import { runUnitOn } from '../utils/helpers'

const TEST_NAME = 'bearer-component-tag-scoping'

describe('NO scope provided', () => {
  runUnitOn(TEST_NAME)
})

describe('With provided scenario ID', () => {
  runUnitOn(TEST_NAME, { tagNamePrefix: 'bearer', tagNameSuffix: 'xyz' })
})
