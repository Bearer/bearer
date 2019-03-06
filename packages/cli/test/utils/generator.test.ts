import * as path from 'path'

import { OpenApiSpecGenerator } from '../../src/utils/generators'

describe('generators', () => {
  it('exports Spec Generator', () => {
    expect(OpenApiSpecGenerator).toBeTruthy()
  })
  // TODO: CORE-197

  describe('build', () => {
    it('generates openapi documentation', async () => {
      const generator = new OpenApiSpecGenerator(path.join(__dirname, '__fixtures__/generators'), {
        integrationTitle: 'test',
        integrationUuid: '123-test'
      })
      const build = await generator.build()
      expect(build).toMatchSnapshot()
    })
  })
})
