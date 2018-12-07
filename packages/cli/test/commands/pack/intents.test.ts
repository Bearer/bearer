import { buildIntentDeclaration } from '../../../src/commands/pack/intents'

describe('pack:intents', () => {
  describe('lambdas declaration', () => {
    it('generates a correct output', () => {
      const expectedOutput = `const intent0 = require(\"./dist/my-intent\").default;
module.exports['my-intent'] = intent0.intentType.intent(intent0.action);

const intent1 = require(\"./dist/spongeBobIntent\").default;
module.exports['spongeBobIntent'] = intent1.intentType.intent(intent1.action);
`

      expect(buildIntentDeclaration({ intents: ['my-intent', 'spongeBobIntent'] })).toEqual(expectedOutput)
    })
  })
})
