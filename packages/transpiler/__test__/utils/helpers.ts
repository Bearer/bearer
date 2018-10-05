import fs from 'fs'
import path from 'path'

import { TranpilerOptions } from '../../src/index'
import { BuildUnitFixtureDirectory, UnitFixtureDirectory } from '../utils/location'

import { TranspilerFactory } from './transpiler'

const fixtures = path.join(__dirname, '__fixtures__')

export function runUnitOn(name: string, transpilerOptions: Partial<TranpilerOptions> = {}) {
  const srcDirectory = UnitFixtureDirectory(name)
  const buildDirectory = BuildUnitFixtureDirectory(name)

  beforeAll(() => {
    runTranspiler(transpilerOptions)
  })

  fs.readdirSync(srcDirectory).forEach(file => {
    describe(file.replace(/\.tsx?/, '').replace(/\-/g, ' '), () => {
      it('match snapshot', async done => {
        expect.assertions(1)
        fs.readFile(path.join(buildDirectory, file), 'utf8', (_e, postContent) => {
          expect({
            postContent,
            file
          }).toMatchSnapshot()
          done()
        })
      })
    })
  })
}

export function runTranspiler(transpilerOptions: Partial<TranpilerOptions> = {}) {
  const transpiler = TranspilerFactory({
    ...transpilerOptions,
    ROOT_DIRECTORY: fixtures,
    srcFolder: '../../__fixtures__'
  })

  transpiler.run()
}
