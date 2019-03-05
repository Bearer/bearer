import * as path from 'path'

import prepareConfig from '../../src/utils/prepare-config'

const base = path.join(__dirname, './__fixtures__/prepare-config')

describe('prepareConfig', () => {
  it('works', async () => {
    const config = await prepareConfig(path.join(base, 'auth.json'), '123-ok', path.join(base, 'intents'))
    expect(config).toMatchObject({
      auth: {},
      integration_uuid: '123-ok',
      intents: expect.arrayContaining(['a-wonderful-intent-as-method', 'a-wonderful-intent'])
    })
  })
})
