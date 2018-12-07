import { buildLambdaDefinitions } from '../../../src/commands/pack/intents'

describe('pack:intents', () => {
  describe('lambdas declaration', () => {
    const { config, handlers } = buildLambdaDefinitions(['my-intent', 'spongeBobIntent'])

    it('generates a correct handlers', () => {
      const expectedOutput = `const intent0 = require(\"./dist/my-intent\").default;
module.exports['my-intent'] = intent0.intentType.intent(intent0.action);

const intent1 = require(\"./dist/spongeBobIntent\").default;
module.exports['spongeBobIntent'] = intent1.intentType.intent(intent1.action);
`
      expect(handlers).toEqual(expectedOutput)
    })

    it('build intents bearer config', () => {
      expect(config).toEqual({
        intents: [
          {
            'my-intent': 'index.my-intent'
          },
          {
            spongeBobIntent: 'index.spongeBobIntent'
          }
        ]
      })
    })
  })
})
