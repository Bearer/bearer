import { Component, State, Prop, Listen, Method } from '@bearer/core'

@Component({
  tag: 'bearer-button-popover',
  styleUrl: 'button-popover.scss',
  shadow: true
})
export class BearerButtonPopover {
  @State() visible: boolean = process.env.NODE_ENV === 'development'
  @Prop() opened: boolean
  @Prop() direction: string = 'top'
  @Prop() arrow: boolean = true
  @Prop() header: string
  @Prop() backNav: boolean
  @Prop()
  btnKind:
    | 'primary'
    | 'secondary'
    | 'success'
    | 'danger'
    | 'warning'
    | 'info'
    | 'light'
    | 'dark' =
    'primary'

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
    if (this.opened === false) {
      this.visible = false
    }
  }

  render() {
    return (
      <div class="root">
        <bearer-button onClick={this.toggleDisplay} kind={this.btnKind}>
          <slot name="buttonText" />
        </bearer-button>
        {this.visible && (
          <div
            class={`popover fade show bs-popover-${this.direction} direction-${
              this.direction
            }`}
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
        )}
      </div>
    )
  }
}
