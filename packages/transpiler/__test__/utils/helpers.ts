import * as fs from 'fs'
import * as globby from 'globby'
import * as path from 'path'
import * as ts from 'typescript'

import { TranpilerOptions } from '../../src/index'
import { UnitFixtureDirectory } from '../utils/location'

import { TranspilerFactory } from './transpiler'

const ROOT_DIRECTORY = path.join(__dirname, '..')
const buildSrcFolder = path.join(ROOT_DIRECTORY, '.build/src')

export function cleanBuildFolder() {
  const buildFolder = path.join(__dirname, '..', '.build/src')

  globby.sync(['**/*.tsx', '**/*.ts', '**/*.json'], { cwd: buildFolder }).forEach(file => {
    const filePath = path.join(buildFolder, file)
    if (fs.existsSync(filePath)) {
      fs.unlinkSync(filePath)
    }
  })
}

export function runTransformers(code: string, transformers: ts.TransformerFactory<ts.SourceFile>[]) {
  const sourceFile = ts.createSourceFile('tmp.ts', code, ts.ScriptTarget.Latest)
  const transformed = ts.transform(sourceFile, transformers)

  const printer = ts.createPrinter(
    {
      newLine: ts.NewLineKind.LineFeed
    },
    {
      onEmitNode: transformed.emitNodeWithNotification,
      substituteNode: transformed.substituteNode
    }
  )
  const result = printer.printBundle(ts.createBundle(transformed.transformed))
  transformed.dispose()

  return result
}

export function runUnitOn(name: string, transpilerOptions: Partial<TranpilerOptions> = {}) {
  const srcDirectory = UnitFixtureDirectory(name)

  beforeAll(() => {
    cleanBuildFolder()
    runTranspiler(`__fixtures__/unit/${name}`, transpilerOptions)
  })

  fs.readdirSync(srcDirectory).forEach(file => {
    describe(file.replace(/\.tsx?/, '').replace(/\-/g, ' '), () => {
      it('match snapshot', async done => {
        expect.assertions(1)
        const filePath = path.join(buildSrcFolder, file)
        fs.readFile(filePath, 'utf8', (_e, postContent) => {
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

export function runTranspiler(srcFolder: string, transpilerOptions: Partial<TranpilerOptions> = {}) {
  const transpiler = TranspilerFactory({
    buildFolder: '.build',
    srcFolder,
    ROOT_DIRECTORY,
    ...transpilerOptions
  })

  transpiler.run()
}
