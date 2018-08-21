import { test } from '@oclif/test'
import { expect } from 'fancy-test'

import { ensureBearerStructure } from '../../helpers/setup'

describe('Generate', () => {
  let bearerPath = ensureBearerStructure()

  beforeEach(done => {
    ensureBearerStructure()
    done()
  })

  describe('generate:setup', () => {
    test
      .stdout()
      .command(['generate:setup', '--force', '--path', bearerPath])
      .it('Generate setup files intent', ctx => {
        expect(ctx.stdout).to.contain('Setup components successfully generated!')
      })
  })
})
