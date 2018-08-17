import { test } from '@oclif/test'
import { expect } from 'fancy-test'

describe('login', () => {
  describe('Prompt for token', () => {
    test
      .stdout()
      .command(['link', '123-scenario-id'])
      .it('Display success message', ctx => {
        expect(ctx.stdout).to.contain('Scenario successfully linked')
      })
  })
})
