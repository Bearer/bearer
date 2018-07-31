import fs from 'fs'
import path from 'path'

const fixtures = path.join(__dirname, '__fixtures__')
const preFolder = path.join(fixtures, 'pre')
const postFolder = path.join(fixtures, '../../.build/src/pre')

describe('Transpiler integration test', () => {
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
        fs.readFile(path.join(postFolder, file), 'utf8', (e, postContent) => {
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
