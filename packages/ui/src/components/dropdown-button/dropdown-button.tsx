import { Component, Listen, Method, Prop, State } from '@bearer/core'
import Bearer from '@bearer/core'

@Component({
  tag: 'bearer-dropdown-button',
  styleUrl: 'dropdown-button.scss',
  shadow: true
})
export class BearerDropdownButton {
  @State() visible: boolean = process.env.NODE_ENV === 'development'
  @Prop() opened: boolean
  @Prop() innerListener: string
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
    if (this.innerListener) {
      Bearer.emitter.addListener(this.innerListener, () => {
        this.visible = false
      })
    }
  }

  render() {
    return (
      <div class="root">
        <bearer-button onClick={this.toggleDisplay} kind={this.btnKind}>
          <slot name="buttonText" />
          <span class="symbol">▾</span>
        </bearer-button>
        {this.visible && (
          <div class="dropdown-down">
            <slot />
          </div>
        )}
      </div>
    )
  }
}
