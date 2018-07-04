import Transpiler from '../src/index'
import * as fs from 'fs-extra'
import * as path from 'path'
import * as ts from 'typescript'
import decoratorNames from './utils'

const BUILD_DIRECTORY = path.join(__dirname, '.build')
const SRC_DIRECTORY = path.join(__dirname, 'files')

afterEach(() => {
  fs.remove(BUILD_DIRECTORY)
})

test('invoking transpiler', async () => {
  let options = {
    noEmitOnError: true,
    noImplicitAny: true,
    target: ts.ScriptTarget.ES5,
    module: ts.ModuleKind.CommonJS,
    outDir: path.resolve(BUILD_DIRECTORY)
  }

  let filePaths = [
    path.join(SRC_DIRECTORY, 'exportObject.ts'),
    path.join(SRC_DIRECTORY, 'classComponent.ts')
  ]
  let transpiler = new Transpiler(options)
  await transpiler.run(filePaths)

  expect(
    fs.existsSync(path.join(BUILD_DIRECTORY, 'exportObject.ts'))
  ).toBeTruthy()

  expect(
    fs.existsSync(path.join(BUILD_DIRECTORY, 'classComponent.ts'))
  ).toBeTruthy()
})

test('Adding BEARER_ID and SCENARIO_ID props', async () => {
  let options = {
    noEmitOnError: true,
    noImplicitAny: true,
    target: ts.ScriptTarget.ES5,
    module: ts.ModuleKind.CommonJS,
    outDir: path.resolve(BUILD_DIRECTORY)
  }

  let filePaths = [path.join(SRC_DIRECTORY, 'classComponent.ts')]
  let transpiler = new Transpiler(options)
  await transpiler.run(filePaths)

  const builtFilePath = path.join(BUILD_DIRECTORY, 'classComponent.ts')
  const program = ts.createProgram([builtFilePath], options)
  let sourceFile = program.getSourceFile(builtFilePath)
  expect(decoratorNames(sourceFile)).toEqual(['Component', 'Prop', 'Prop'])
})
