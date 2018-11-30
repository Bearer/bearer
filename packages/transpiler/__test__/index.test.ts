import * as fs from 'fs'
import * as path from 'path'

import { runTranspiler } from './utils/helpers'

const fixtures = path.join(__dirname, '__fixtures__')
const preFolder = path.join(fixtures, 'pre')
const postFolder = path.join(fixtures, '../../.build/src/pre')

describe('Transpiler integration test', () => {
  beforeAll(done => {
    runTranspiler()
    done()
  })
  fs.readdirSync(preFolder).forEach(file => {
    describe(file, () => {
      it(`creates a file `, done => {
        expect.assertions(1)
        fs.exists(path.join(postFolder, file), exists => {
          expect(exists).toBe(true)
          done()
        })
      })

      it('match expectation', done => {
        expect.assertions(1)
        fs.readFile(path.join(postFolder, file), 'utf8', (_e, postContent) => {
          expect({
            postContent,
            file
          }).toMatchSnapshot()
          done()
        })
      })

      it('has correct manifest generated', done => {
        fs.readFile(path.join(postFolder, '../bearer-manifest.json'), 'utf8', (_e, postContent) => {
          const manifest = JSON.parse(postContent).manifest.components.map(({ outputs, inputs, classname }) => ({
            classname,
            inputs,
            outputs
          }))
          expect({
            manifest,
            file
          }).toMatchSnapshot()
          done()
        })
      })
    })
  })
})
