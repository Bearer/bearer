import { test } from '@oclif/test'
import { expect } from 'fancy-test'
import * as sinon from 'sinon'

import * as setup from '../../src/utlis/setupConfig'

describe('Link', () => {
  describe('Update', () => {
    const update = sinon.spy()
    const stub = sinon.stub(setup, 'default')
    stub.returns({
      setScenarioConfig: update
    })

    test
      .stdout()
      .command(['link', '123-scenario-id'])
      .it('Display success message', ctx => {
        expect(update.args[0][0]).to.include({
          orgId: '123',
          scenarioId: 'scenario-id'
        })
        stub.restore()
        expect(ctx.stdout).to.contain('Scenario successfully linked')
      })
  })
})
