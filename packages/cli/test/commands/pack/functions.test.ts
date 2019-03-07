import { buildLambdaDefinitions } from '../../../src/commands/pack/functions'

describe('pack:functions', () => {
  describe('lambdas declaration', () => {
    const { config, handlers } = buildLambdaDefinitions(['my-intent', 'spongeBobFunction'])

    it('generates a correct handlers', () => {
      expect(handlers).toMatchSnapshot()
    })

    it('build functions bearer config', () => {
      expect(config).toEqual({
        functions: [
          {
            'my-intent': 'index.my-intent'
          },
          {
            spongeBobFunction: 'index.spongeBobFunction'
          }
        ]
      })
    })
  })
})
