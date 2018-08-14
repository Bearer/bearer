import { expect, test } from '@oclif/test'

describe('new', () => {
  test
    .stdout()
    .command(['new', '.bearer/SpongeBobScenario', '-a', 'oauth2'])
    .it('generates a scenario without any prompt', ctx => {
      expect(ctx.stdout).to.contain('Scenario initialized, name: SpongeBobScenario, authentication type: OAuth2')
    })
})
