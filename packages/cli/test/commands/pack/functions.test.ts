import { buildLambdaDefinitions } from '../../../src/commands/pack/functions'

describe('pack:functions', () => {
  describe('lambdas declaration', () => {
    const { config } = buildLambdaDefinitions(['my-function', 'spongeBobFunction'])

    it('build functions bearer config', () => {
      expect(config).toEqual({
        functions: [
          {
            'my-function': 'index.my-function'
          },
          {
            spongeBobFunction: 'index.spongeBobFunction'
          }
        ]
      })
    })
  })
})
