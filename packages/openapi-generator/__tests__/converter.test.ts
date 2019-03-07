import fs from 'fs'
import path from 'path'
import ts from 'typescript'

import { functionTypesToSchemaConverter as converter } from '../src/index'
const FUNCTIONS_DIR = path.join(__dirname, '__fixtures__', 'integration', 'functions')

const functions = fs.readdirSync(FUNCTIONS_DIR).map(func => [func])

describe('#functionTypeToSchemaConverter', () => {
  test.each(functions)('converts types to schemas for %s', file => {
    expect(
      converter(path.join(FUNCTIONS_DIR, file), {
        target: ts.ScriptTarget.ES5,
        module: ts.ModuleKind.CommonJS
      })
    ).toMatchSnapshot()
  })
})
