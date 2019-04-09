import * as path from 'path'

import prepareConfig from '../../src/utils/prepare-config'

const base = path.join(__dirname, './__fixtures__/prepare-config')

describe('prepareConfig', () => {
  it('works', async () => {
    const config = await prepareConfig(path.join(base, 'auth.json'), '123-ok', path.join(base, 'functions'))
    expect(config).toMatchObject({
      auth: {},
      buid: '123-ok',
      functions: expect.arrayContaining(['a-wonderful-function-as-method', 'a-wonderful-function'])
    })
  })
})
