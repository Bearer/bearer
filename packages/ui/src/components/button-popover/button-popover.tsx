import { Component, State, Prop, Listen, Method } from '@bearer/core'

@Component({
  tag: 'bearer-button-popover',
  styleUrl: 'button-popover.scss',
  shadow: true
})
export class BearerButtonPopover {
  @State() visible: boolean = false

  @Prop() opened: boolean
  @Prop() direction: string = 'top'
  @Prop() arrow: boolean = true
  @Prop() header: string
  @Prop() backNav: boolean
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
  }

  render() {
    return (
      <div class="root">
        <bearer-button {...this.btnProps} onClick={this.toggleDisplay} />

        <div
          class={`popover fade show bs-popover-${this.direction} direction-${this.direction} ${!this.visible &&
            'hidden'}`}
        >
          <h3 class="popover-header">
            {this.backNav && <bearer-navigator-back class="header-arrow" />}
            <span class="header">{this.header}</span>
          </h3>
          <div class="popover-body">
            <slot />
          </div>
          {this.arrow && <div class="arrow" />}
        </div>
      </div>
    )
  }
}
