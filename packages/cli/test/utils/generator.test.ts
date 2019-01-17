import * as path from 'path'

import { IOpenApiSpec, OpenApiSpecGenerator } from '../../src/utils/generators'

describe('generators', () => {
  it('exports Spec Generator', () => {
    expect(OpenApiSpecGenerator).toBeTruthy()
  })
  // TODO: CORE-197
  describe.skip('params definitions', () => {
    let result: IOpenApiSpec

    beforeAll(async () => {
      const generator = new OpenApiSpecGenerator(path.join(__dirname, '__fixtures__/generators'), {
        scenarioTitle: 'test',
        scenarioUuid: '123-test'
      })
      result = await generator.build()
    })

    it('has object literal params required', () => {
      const paramsSchema = result.paths[`/123-test/object-literal-type`].post.parameters.find(
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

    // TODO: CORE-197
    it('has object literal params optional', () => {
      const optionalParam = result.paths[`/123-test/object-literal-type`].post.parameters.find(
        p => p.name === 'optional'
      )
      expect(optionalParam).toMatchObject({
        description: 'optional',
        in: 'query',
        name: 'optional',
        required: false,
        schema: { type: 'number' }
      })
    })

    // TODO: CORE-197
    it('has aliased type params', () => {
      const paramsSchema = result.paths[`/123-test/type-alias`].post.parameters.find(p => p.name === 'aliasedParams')
      expect(paramsSchema).toBeTruthy()
      expect(paramsSchema).toMatchObject({
        description: 'aliasedParams',
        in: 'query',
        name: 'aliasedParams',
        required: true,
        schema: { type: 'string' }
      })
    })

    // TODO: CORE-197
    it('returns undefined params and return', () => {
      const paramsSchema = result.paths[`/123-test/undefined-params-and-return`].post.parameters
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
