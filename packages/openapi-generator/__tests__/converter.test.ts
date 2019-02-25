import fs from 'fs'
import path from 'path'
import ts from 'typescript'

import { intentTypesToSchemaConverter as converter } from '../src/index'
const INTENTS_DIR = path.join(__dirname, '__fixtures__', 'scenario', 'intents')

const intents = fs.readdirSync(INTENTS_DIR).map(intent => [intent])

describe('#intentTypeToSchemaConverter', () => {
  test.each(intents)('converts types to schemas for %s', file => {
    expect(
      converter(path.join(INTENTS_DIR, file), {
        target: ts.ScriptTarget.ES5,
        module: ts.ModuleKind.CommonJS
      })
    ).toMatchSnapshot()
  })
})
