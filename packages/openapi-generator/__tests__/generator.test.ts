import fs from 'fs'
import path from 'path'

import generator from '../src'

const INTENTS_DIR = path.resolve(path.join(__dirname, '__fixtures__', 'scenario', 'intents'))

const intents = fs.readdirSync(INTENTS_DIR).map(intent => intent.replace(/\.ts$/, ''))

describe('#generator', () => {
  test('generates the valid openapi 3.0 document', () => {
    const integrationUuid = '4l1c3-happy-goats'
    const integrationName = 'happy-goats'
    const document = generator({ intents, integrationUuid, integrationName, intentsDir: INTENTS_DIR })
    expect(document).toMatchSnapshot()
  })
})
