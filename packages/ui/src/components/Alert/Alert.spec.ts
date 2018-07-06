import { render } from '@stencil/core/testing'
import { Alert } from './Alert'

describe('alert', () => {
  it('should build', () => {
    expect(new Alert()).toBeTruthy()
  })

  describe('rendering', () => {
    let element
    beforeEach(async () => {
      element = await render({
        components: [Alert],
        html: '<bearer-alert>Sponge bob</bearer-alert>'
      })
    })

    it('works without parameters', () => {
      expect(element.textContent.trim()).toEqual('Sponge bob')
    })
  })
})
