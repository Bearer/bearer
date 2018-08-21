import { test } from '@oclif/test'
import { expect } from 'fancy-test'

import { ensureBearerStructure } from '../../helpers/setup'

describe('Generate', () => {
  let bearerPath = ensureBearerStructure()
  beforeEach(done => {
    ensureBearerStructure()
    done()
  })
  describe('generate:intent', () => {
    test
      .stdout()
      .command(['generate:intent', 'FetchDataIntent', '-t', 'fetch', '--path', bearerPath])
      .it('Fetch intent', ctx => {
        expect(ctx.stdout).to.contain('Generated intent: name: FetchDataIntent type: FetchData')
      })

    test
      .stdout()
      .command(['generate:intent', 'SaveDataIntent', '-t', 'save', '--path', bearerPath])
      .it('Save Intent', ctx => {
        expect(ctx.stdout).to.contain('Generated intent: name: SaveDataIntent type: SaveState')
      })

    test
      .stdout()
      .command(['generate:intent', 'RetrieveDataIntent', '-t', 'retrieve', '--path', bearerPath])
      .it('Retrieve Intent', ctx => {
        expect(ctx.stdout).to.contain('Generated intent: name: RetrieveDataIntent type: RetrieveState')
      })
  })
})
