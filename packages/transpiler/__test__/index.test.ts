import fs from 'fs'
import Transpiler from '../src'
import path from 'path'
import { idText } from 'typescript'

const fixtures = path.join(__dirname, '__fixtures__')
const preFolder = path.join(fixtures, 'pre')
const postFolder = path.join(fixtures, 'post')
const postExpectFolder = path.join(fixtures, 'post-expectations')

const transpiler = new Transpiler(
  __dirname + '/__fixtures__/pre/',
  false,
  '../post/'
)

describe('Transpiler integration test', () => {
  beforeAll(() => {
    fs.readdirSync(postFolder).forEach(file => {
      if (file !== 'tsconfig.json') {
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
          fs.readFile(
            path.join(postExpectFolder, file),
            'utf8',
            (e, postContentExpected) => {
              expect(postContent).toMatch(postContentExpected)
              done()
            }
          )
        })
      })
    })
  })
})
// test('invoking transpiler', async () => {
//   let transpiler = new Transpiler(SRC_DIRECTORY)
//   await transpiler.run()
//   expect(
//     fs.existsSync(path.join(BUILD_DIRECTORY, 'exportObject.ts'))
//   ).toBeTruthy()
//   expect(
//     fs.existsSync(path.join(BUILD_DIRECTORY, 'classComponent.ts'))
//   ).toBeTruthy()
// })

// test('Adding BEARER_ID prop', async () => {
//   pending('circular calls')
//   let transpiler = new Transpiler(SRC_DIRECTORY)
//   await transpiler.run()
//   const builtFilePath = path.join(BUILD_DIRECTORY, 'classComponent.ts')
// })
