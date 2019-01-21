import * as path from 'path'

import { IOpenApiSpec, OpenApiSpecGenerator } from '../../src/utils/generators'

describe('generators', () => {
  it('exports Spec Generator', () => {
    expect(OpenApiSpecGenerator).toBeTruthy()
  })
  // TODO: CORE-197
  const expectedEndpointLength = 3

  let result: IOpenApiSpec

  beforeAll(async () => {
    const generator = new OpenApiSpecGenerator(path.join(__dirname, '__fixtures__/generators'), {
      scenarioTitle: 'test',
      scenarioUuid: '123-test'
    })
    result = await generator.build()
  })

  describe('responses', () => {
    it('generate as many entries as intents', () => {
      expect(Object.keys(result.paths)).toHaveLength(expectedEndpointLength)
    })

    it('had default 200, 401 and 403 statuses', () => {
      expect.assertions(expectedEndpointLength * 3)

      Object.keys(result.paths).map(p => {
        expect(result.paths[p].post.responses['200']).toBeTruthy()
        expect(result.paths[p].post.responses['403']).toBeTruthy()
        expect(result.paths[p].post.responses['403']).toBeTruthy()
      })
    })
  })

  describe.skip('response payload', () => {
    function responseContentObject(intent: string) {
      return result.paths[intent].post.requestBody.content['application/json'].schema
    }

    describe('typing not provided', () => {
      it('return empty input', () => {
        expect(responseContentObject('/123-test/undefined-params-and-return')).toMatchObject({ properties: {} })
      })
    })

    describe('typing inlined', () => {
      it('extract from inline', () => {
        expect(responseContentObject('/123-test/object-literal-type')).toMatchObject({
          type: 'object',
          properties: {
            inlineParam: {
              type: 'string'
            },
            stringEnum: {
              type: 'string',
              enum: ['none', 'all', 'every']
            }
          }
        })
      })
    })

    describe('typing aliased', () => {
      it('extract from alias', () => {
        expect(responseContentObject('/123-test/type-alias')).toMatchObject({
          type: 'object',
          properties: {
            foo: {
              type: 'array',
              required: true,
              items: {
                type: 'string'
              }
            },
            anObject: {
              required: true,
              type: 'object',
              properties: {
                values: {
                  type: 'array',

                  items: {
                    type: 'number'
                  }
                }
              }
            }
          }
        })
      })
    })
  })

  describe.skip('requests', () => {
    function randomResponse() {
      return result.paths['/123-test/object-literal-type'].post
    }

    describe('parameters', () => {
      it('expect an Authorization header', () => {
        expect(randomResponse().parameters[0]).toMatchObject({
          description: 'API Key',
          in: 'header',
          name: 'authorization',
          required: true,
          schema: { type: 'string' }
        })
      })
    })

    describe('post requestBody from intents', () => {
      function requestBodySchema(intent: string) {
        return result.paths[intent].post.requestBody.content['application/json'].schema
      }

      describe('typing not provided', () => {
        it('return empty input', () => {
          expect(requestBodySchema('/123-test/undefined-params-and-return')).toMatchObject({ properties: {} })
        })
      })

      describe('typing inlined', () => {
        it('extract from inline', () => {
          expect(requestBodySchema('/123-test/object-literal-type')).toMatchObject({
            type: 'object',
            properties: {
              inlineParam: {
                type: 'string'
              },
              stringEnum: {
                type: 'string',
                enum: ['none', 'all', 'every']
              }
            }
          })
        })
      })

      describe('typing aliased', () => {
        it('extract from alias', () => {
          expect(requestBodySchema('/123-test/type-alias')).toMatchObject({
            type: 'object',
            properties: {
              aliasParam: {
                type: 'string'
              },
              stringEnum: {
                type: 'string',
                enum: ['none', 'all', 'every']
              }
            }
          })
        })
      })
    })
  })
})
