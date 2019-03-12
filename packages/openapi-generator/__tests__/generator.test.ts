import fs from 'fs'
import path from 'path'

import generator from '../src'

const FUNCTIONS_DIR = path.resolve(path.join(__dirname, '__fixtures__', 'integration', 'functions'))

const functions = fs.readdirSync(FUNCTIONS_DIR).map(func => func.replace(/\.ts$/, ''))

describe('#generator', () => {
  test('generates the valid openapi 3.0 document', () => {
    const integrationUuid = '4l1c3-happy-goats'
    const integrationName = 'happy-goats'
    const document = generator({ functions, integrationUuid, integrationName, functionsDir: FUNCTIONS_DIR })
    expect(document).toMatchSnapshot()
  })
})
