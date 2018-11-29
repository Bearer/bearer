import * as path from 'path'

import { IOpenApiSpec, OpenApiSpecGenerator } from '../../src/utils/generators'

describe('generators', () => {
  it('exports Spec Generator', () => {
    expect(OpenApiSpecGenerator).toBeTruthy()
  })
  describe('params definitions', () => {
    const generator = new OpenApiSpecGenerator(path.join(__dirname, '__fixtures__/generators'), {
      scenarioTitle: 'test',
      scenarioUuid: '123-test'
    })
    let result: IOpenApiSpec

    beforeAll(async () => {
      result = await generator.build()
    })

    it.only('has object literal params', () => {
      const paramsSchema = result.paths[`/123-test/ObjectLiteralParams`].post.parameters.find(
        p => p.name === 'inlineParam'
      )
      expect(paramsSchema).toBeTruthy()
      expect(paramsSchema).toMatchObject({
        description: 'inlineParam',
        in: 'query',
        name: 'inlineParam',
        required: true,
        schema: { type: 'string' }
      })
    })

    it('has aliased type params', () => {
      const paramsSchema = result.paths[`/123-test/TypeAliasParams`].post.parameters.find(
        p => p.name === 'aliasedParams'
      )
      expect(paramsSchema).toBeTruthy()
      expect(paramsSchema).toMatchObject({
        description: 'aliasedParams',
        in: 'query',
        name: 'aliasedParams',
        required: true,
        schema: { type: 'string' }
      })
    })
  })
})
