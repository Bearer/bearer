import { Component, Listen, Method, Prop, Element } from '@bearer/core'

@Component({
  tag: 'bearer-dropdown-button'
})
export class BearerDropdownButton {
  @Element() el: HTMLElement
  @Prop() opened: boolean = false
  @Prop() innerListener: string
  @Prop() btnProps: JSXElements.BearerButtonAttributes = {}
  @Prop() backNav: boolean = false
  @Listen('body:click')
  clickOutsideHandler() {
    this.toggle(false)
  }

  @Listen('click')
  clickInsideHandler(ev) {
    ev.stopImmediatePropagation()
  }

  @Method()
  toggle(opened: boolean) {
    this.popover.toggle(opened)
  }

  close = () => {
    this.toggle(false)
  }

  componentDidLoad() {
    if (this.innerListener) {
      this.el.addEventListener(this.innerListener, this.close)
    }
  }

  componentDidUnload() {
    if (this.innerListener) {
      this.el.removeEventListener(this.innerListener, this.close)
    }
  }

  private get popover(): HTMLBearerDropdownButtonElement {
    return this.el.querySelector<HTMLBearerDropdownButtonElement>('bearer-button-popover')
  }

  render() {
    const { content, ...rest } = this.btnProps
    const btnProps: JSXElements.BearerButtonAttributes = {
      kind: 'primary',
      ...rest
    }
    return (
      <bearer-button-popover
        btnProps={btnProps}
        direction="bottom"
        aligned="left"
        opened={this.opened}
        backNav={this.backNav}
      >
        <span slot="btn-content">
          <slot name="dropdown-btn-content" />
          <span class="symbol" style={{ paddingLeft: '10px' }}>
            ▾
          </span>
        </span>
        <slot />
      </bearer-button-popover>
    )
  }
}
