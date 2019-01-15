import Bearer, { Component, Listen, Method, Prop, State } from '@bearer/core'

@Component({
  tag: 'bearer-dropdown-button',
  styleUrl: 'dropdown-button.scss',
  shadow: true
})
export class BearerDropdownButton {
  @State() visible: boolean = false
  @Prop() opened: boolean = false
  @Prop() innerListener: string
  @Prop() btnProps: JSXElements.BearerButtonAttributes = {}

  toggleDisplay = e => {
    e.preventDefault()
    this.visible = !this.visible
  }

  @Listen('body:click')
  clickOutsideHandler() {
    this.visible = false
  }

  @Listen('click')
  clickInsideHandler(ev) {
    ev.stopImmediatePropagation()
  }

  @Method()
  toggle(opened: boolean) {
    this.visible = opened
  }

  componentDidLoad() {
    this.visible = this.opened
    if (this.innerListener) {
      Bearer.emitter.addListener(this.innerListener, () => {
        this.visible = false
      })
    }
  }

  render() {
    const { content, ...rest } = this.btnProps
    const btnProps: JSXElements.BearerButtonAttributes = {
      kind: 'primary',
      ...rest
    }
    return (
      <bearer-button-popover btnProps={btnProps} direction="bottom" aligned="left">
        <slot />
        <span slot="btn-content">
          <slot name="btn-content" />
          <span class="symbol">▾</span>
        </span>
      </bearer-button-popover>
    )
  }
}
