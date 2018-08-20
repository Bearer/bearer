import { test } from '@oclif/test'
import { expect } from 'fancy-test'
import * as sinon from 'sinon'

import * as setup from '../../../src/utils/setupConfig'

describe('Generate', () => {
  let stub: any
  let update: any

  afterEach(() => {
    ;(setup.default as any).restore()
    stub.restore()
  })

  describe('generate:setup', () => {
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
      .command(['generate:setup', '--force'])
      .it('Generate setup files intent', ctx => {
        expect(ctx.stdout).to.contain('Setup components successfully generated!')
      })
  })
})
