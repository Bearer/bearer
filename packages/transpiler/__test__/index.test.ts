import fs from 'fs'
import Transpiler, { TranpilerOptions } from '../src'
import path from 'path'

const fixtures = path.join(__dirname, '__fixtures__')
const preFolder = path.join(fixtures, 'pre')
const postFolder = path.join(fixtures, '../../.build/src')

const options: TranpilerOptions = {
  ROOT_DIRECTORY: fixtures,
  watchFiles: false,
  buildFolder: '../../.build/',
  srcFolder: 'pre'
}
const transpiler = new Transpiler(options)

describe('Transpiler integration test', () => {
  beforeAll(() => {
    process.env.BEARER_SCENARIO_ID = 'SPONGE_BOB'
    console.log('[BEARER]', 'postFolder', postFolder)
    fs.readdirSync(postFolder).forEach(file => {
      if (file !== 'tsconfig.json' && file !== '.keep') {
        fs.unlinkSync(path.join(postFolder, file))
      }
    })
    transpiler.run()
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
