import { runTransformers } from '../utils/helpers'
import { UnitFixtureDirectory } from '../utils/location'

import Metadata from '../../src/metadata'
import GatherMetadata from '../../src/transformers/gather-metadata'
import InputDecorator from '../../src/transformers/input-decorator'
import * as fs from 'fs'
import * as path from 'path'
const TEST_NAME = 'input-decorator'

const metadata = new Metadata("bearer", "test")

describe('Input decorator units', () => {
  it("generates stuff", (done) => {
    const srcDirectory = UnitFixtureDirectory(TEST_NAME)

    fs.readdirSync(srcDirectory).forEach(file => {
      const filePath = path.join(srcDirectory, file)

      fs.readFile(filePath, 'utf8', (_e, postContent) => {
        const result = runTransformers(postContent, [GatherMetadata({ metadata }), InputDecorator({ metadata })])
        expect(result).toMatchSnapshot()
        done()
      })
    })
  })
})
