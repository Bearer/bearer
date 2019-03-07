import { buildLambdaDefinitions } from '../../../src/commands/pack/functions'

describe('pack:functions', () => {
  describe('lambdas declaration', () => {
    const { config, handlers } = buildLambdaDefinitions(['my-function', 'spongeBobFunction'])

    it('generates a correct handlers', () => {
      expect(handlers).toMatchSnapshot()
    })

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
