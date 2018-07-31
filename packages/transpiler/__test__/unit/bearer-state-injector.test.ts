import fs from 'fs'
import { UnitFixtureDirectory, BuildUnitFixtureDirectory } from '../utils/location'
import path from 'path'

const TEST_NAME = 'bearer-state-injector'
const srcDirectory = UnitFixtureDirectory(TEST_NAME)
const buildDirectory = BuildUnitFixtureDirectory(TEST_NAME)

describe('Bearer State Injector', () => {
  fs.readdirSync(srcDirectory).forEach(file => {
    describe(file.replace('.ts', '').replace(/\-/g, ' '), () => {
      it('match snapshot', async done => {
        expect.assertions(1)
        fs.readFile(path.join(buildDirectory, file), 'utf8', (e, postContent) => {
          expect({
            postContent,
            file
          }).toMatchSnapshot()
          done()
        })
      })
    })
  })
})
