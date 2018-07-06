import { render } from '@stencil/core/testing'
import { BearerBadge } from './Badge'

describe('badge', () => {
  it('should build', () => {
    expect(new BearerBadge()).toBeTruthy()
  })

  describe('rendering', () => {
    let element
    beforeEach(async () => {
      element = await render({
        components: [BearerBadge],
        html: '<bearer-badge>Sponge Bobd</bearer-badge>'
      })
    })

    it('should work without parameters', () => {
      expect(element.textContent.trim()).toEqual('Sponge Bobd')
    })
  })
})
