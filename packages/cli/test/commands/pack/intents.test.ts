import { buildLambdaDefinitions } from '../../../src/commands/pack/intents'

describe('pack:intents', () => {
  describe('lambdas declaration', () => {
    const { config, handlers } = buildLambdaDefinitions(['my-intent', 'spongeBobFunction'])

    it('generates a correct handlers', () => {
      expect(handlers).toMatchSnapshot()
    })

    it('build intents bearer config', () => {
      expect(config).toEqual({
        intents: [
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
