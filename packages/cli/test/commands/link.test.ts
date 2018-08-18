import { test } from '@oclif/test'
import { expect } from 'fancy-test'
import * as sinon from 'sinon'

import * as setup from '../../src/utils/setupConfig'

describe('Link', () => {
  let stub: any
  let update: any

  afterEach(() => {
    ;(setup.default as any).restore()
    stub.restore()
  })

  describe('Update', () => {
    beforeEach(() => {
      update = sinon.spy()
      stub = sinon.stub(setup, 'default')
      stub.returns({
        setScenarioConfig: update,
        isScenarioLocation: true
      })
    })

    test
      .stdout()
      .command(['link', '123-scenario-id'])
      .it('Display success message', ctx => {
        expect(update.args[0][0]).to.include({
          orgId: '123',
          scenarioId: 'scenario-id'
        })
        expect(ctx.stdout).to.contain('Scenario successfully linked')
      })
  })

  describe('Does not update', () => {
    beforeEach(() => {
      update = sinon.spy()
      stub = sinon.stub(setup, 'default')
      stub.returns({
        setScenarioConfig: update,
        isScenarioLocation: false
      })
    })

    test
      .stdout()
      .stderr()
      .command(['link', '123-scenario-id'])
      .exit(2)
      .it('exits with error')
  })
})
