import { runTransformers } from '../utils/helpers'
import { UnitFixtureDirectory } from '../utils/location'

import Metadata from '../../src/metadata'
import GatherMetadata from '../../src/transformers/gather-metadata'
import OutputDecorator from '../../src/transformers/output-decorator'
import * as fs from 'fs'
import * as path from 'path'
const TEST_NAME = 'input-output-intent-arguments'

const metadata = new Metadata('bearer', 'test')

describe('input and output arguments', () => {
  it('generates the refIds properly', done => {
    const srcDirectory = UnitFixtureDirectory(TEST_NAME)

    fs.readdirSync(srcDirectory).forEach(file => {
      const filePath = path.join(srcDirectory, file)

      fs.readFile(filePath, 'utf8', (_e, postContent) => {
        const result = runTransformers(postContent, [GatherMetadata({ metadata }), OutputDecorator({ metadata })])
        expect(result).toMatchSnapshot()
        done()
      })
    })
  })
})
