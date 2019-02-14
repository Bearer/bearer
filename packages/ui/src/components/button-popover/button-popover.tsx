import { Component, Event, EventEmitter, Method, Prop, State, Watch } from '@bearer/core'
import { BKind } from '../Button/Button'

export type TAlignement = 'left' | 'right'
export type TDirection = 'left' | 'right' | 'top' | 'bottom'

@Component({
  tag: 'bearer-button-popover',
  styleUrl: 'button-popover.scss',
  shadow: true
})
export class BearerButtonPopover {
  @Event()
  visibilityChange: EventEmitter

  @Prop({ reflectToAttr: true }) kind: BKind = 'embed'
  @Prop({ reflectToAttr: true, mutable: true }) opened: boolean
  @Prop({ reflectToAttr: true }) direction: TDirection = 'right'
  @Prop({ reflectToAttr: true }) aligned: TAlignement = 'left'

  @Prop() header: string
  @Prop() backNav: boolean = false
  @Prop() btnProps: JSXElements.BearerButtonAttributes = {}

  @State() _visible: boolean = false

  toggleDisplay = e => {
    e.preventDefault()

    console.debug('[BEARER]', 'Button popover: toggleDisplay', !this.visible)
    this.visible = !this.visible
  }

  set visible(newValue: boolean) {
    if (newValue !== null && this._visible !== newValue) {
      console.debug('[BEARER]', 'Button popover: visibilityChangeHandler', newValue)
      this._visible = newValue
      this.visibilityChange.emit({ visible: this._visible })
    }

    this.opened = newValue
  }

  get visible(): boolean {
    return this._visible
  }

  clickInsidePopover(ev) {
    ev.stopImmediatePropagation()
  }

  @Method()
  toggle(newValue: boolean) {
    // Set visibility to toggle param
    // or inverse the current one.
    this.visible = typeof newValue !== 'undefined' ? newValue : !this.visible
  }

  @Watch('opened')
  watchOpened(newValue: boolean) {
    // Opened shall be set (true or false)
    if (newValue === true || newValue === false) {
      this.visible = newValue
    }
  }

  componentDidLoad() {
    if (this.opened !== null) {
      this.visible = this._visible
      this.visible = typeof this.opened !== 'undefined' ? this.opened : this._visible
    }
  }

  render() {
    return [
      <bearer-button kind={this.kind} {...this.btnProps} onClick={this.toggleDisplay}>
        <slot name="btn-content" />
      </bearer-button>,

      <div
        onClick={this.clickInsidePopover}
        class={`popover direction-${this.direction} ${!this.visible && 'hidden'} aligned-${this.aligned}`}
      >
        {(this.backNav || this.header) && (
          <h3 class="popover-header">
            {this.backNav && <bearer-navigator-back class="header-arrow" />}
            <span class="header">{this.header}</span>
          </h3>
        )}
        <div class="popover-body">
          <slot />
        </div>
      </div>
    ]
  }
}
