import { Component, Event, EventEmitter, Listen, Method, Prop, State } from '@bearer/core'
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

  @Prop({ reflectToAttr: true }) kind: BKind = 'action'
  @Prop({ reflectToAttr: true }) opened: boolean
  @Prop({ reflectToAttr: true }) direction: TDirection = 'right'
  @Prop({ reflectToAttr: true }) aligned: TAlignement = 'left'

  @Prop() header: string
  @Prop() backNav: boolean = true
  @Prop() btnProps: JSXElements.BearerButtonAttributes = {}

  @State()
  _visible: boolean = false

  toggleDisplay = e => {
    e.preventDefault()
    console.debug('[BEARER]', 'Button popover: toggleDisplay', !this.visible)
    this.visible = !this.visible
  }

  set visible(newValue: boolean) {
    if (this._visible !== newValue) {
      console.debug('[BEARER]', 'Button popover: visibilityChangeHandler', newValue)
      this._visible = newValue
      this.visibilityChange.emit({ visible: this._visible })
    }
  }

  get visible(): boolean {
    return this._visible
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
    return [
      <bearer-button kind={this.kind} {...this.btnProps} onClick={this.toggleDisplay}>
        <slot name="btn-content" />
      </bearer-button>,

      <div class={`popover direction-${this.direction} ${!this.visible && 'hidden'} aligned-${this.aligned}`}>
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
