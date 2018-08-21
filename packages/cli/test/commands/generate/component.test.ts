import { test } from '@oclif/test'
import { expect } from 'fancy-test'

import { ensureBearerStructure } from '../../helpers/setup'

describe('Generate', () => {
  let bearerPath = ensureBearerStructure()
  beforeEach(done => {
    ensureBearerStructure()
    done()
  })
  describe('generate:component', () => {
    test
      .stdout()
      .command(['generate:component', 'blankComponent', '-t', 'blank', '--path', bearerPath])
      .it('Blank component', ctx => {
        expect(ctx.stdout).to.contain('Component generated')
      })

    test
      .stdout()
      .command(['generate:component', 'collectionComponent', '-t', 'collection', '--path', bearerPath])
      .it('Collection Intent', ctx => {
        expect(ctx.stdout).to.contain('Component generated')
      })
    test
      .stdout()
      .command(['generate:component', 'rootComponent', '-t', 'root', '--path', bearerPath])
      .it('Root Intent', ctx => {
        expect(ctx.stdout).to.contain('Component generated')
      })
  })
})
