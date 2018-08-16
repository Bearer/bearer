import { expect, test } from '@oclif/test'

describe('new command', () => {
  it('generates files and print help', () => {
    test
      .stdout()
      .command(['new', '.bearer/SpongeBobScenario', '-a', 'oauth2', '--skipInstall'])
      .it('generates a scenario without any prompt', ctx => {
        console.log('out', ctx.stdout)
        expect(ctx.stdout).to.contain(`Generate files:`)
        expect(ctx.stdout).to.contain(`create: .bearer/SpongeBobScenario/LICENSE`)
        expect(ctx.stdout).to.contain(
          `Scenario initialized, name: .bearer/SpongeBobScenario, authentication type: OAuth2`
        )
        expect(ctx.stdout).to.contain(`What's next?`)
      })
  })
})
