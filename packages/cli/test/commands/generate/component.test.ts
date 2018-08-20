import { test } from '@oclif/test'
import { expect } from 'fancy-test'
import * as sinon from 'sinon'

import * as setup from '../../../src/utils/setupConfig'

describe('Generare', () => {
  let stub: any
  let update: any

  afterEach(() => {
    ;(setup.default as any).restore()
    stub.restore()
  })

  describe('generate:component', () => {
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
      .command(['generate:component', 'blankComponent', '-t', 'blank'])
      .it('Blank component', ctx => {
        expect(ctx.stdout).to.contain('Generated component: name: blankComponent | type: blank')
      })

    test
      .stdout()
      .command(['generate:component', 'collectionComponent', '-t', 'collection'])
      .it('Collection Intent', ctx => {
        expect(ctx.stdout).to.contain('Generated component: name: collectionComponent | type: collection')
      })
    test
      .stdout()
      .command(['generate:component', 'rootComponent', '-t', 'root'])
      .it('Root Intent', ctx => {
        expect(ctx.stdout).to.contain('Generated component: name: rootComponent | type: root')
      })
  })
})
