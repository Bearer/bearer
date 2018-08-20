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

  describe('generate:intent', () => {
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
      .command(['generate:intent', 'FetchDataIntent', '-t', 'fetch'])
      .it('Fetch intent', ctx => {
        expect(ctx.stdout).to.contain('Generated intent: name: FetchDataIntent type: FetchData')
      })

    test
      .stdout()
      .command(['generate:intent', 'SaveDataIntent', '-t', 'save'])
      .it('Save Intent', ctx => {
        expect(ctx.stdout).to.contain('Generated intent: name: SaveDataIntent type: SaveState')
      })
    test
      .stdout()
      .command(['generate:intent', 'RetrieveDataIntent', '-t', 'retrieve'])
      .it('Retrieve Intent', ctx => {
        expect(ctx.stdout).to.contain('Generated intent: name: RetrieveDataIntent type: RetrieveState')
      })
  })
})
