import { render } from '@stencil/core/testing'
import { Button } from './Button'

describe('button', () => {
  it('should build', () => {
    expect(new Button()).toBeTruthy()
  })

  describe('rendering', () => {
    let element
    beforeEach(async () => {
      element = await render({
        components: [Button],
        html: '<bearer-button>Sponge bob</bearer-button>'
      })
    })

    it('should work without parameters', () => {
      expect(element.textContent.trim()).toEqual('Sponge bob')
    })
  })
})
