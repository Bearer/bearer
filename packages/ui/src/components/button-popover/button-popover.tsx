import { Component, Event, EventEmitter, Listen, Method, Prop, State } from '@bearer/core'
import { BKind } from '../Button/Button'

export type TAlignement = 'left' | 'right'

@Component({
  tag: 'bearer-button-popover',
  styleUrl: 'button-popover.scss',
  shadow: true
})
export class BearerButtonPopover {
  @State()
  _visible: boolean = false
  @Prop({ reflectToAttr: true }) kind: BKind = 'action'
  @Event()
  visibilityChange: EventEmitter
  @Prop({ reflectToAttr: true })
  opened: boolean
  @Prop({ reflectToAttr: true })
  direction: string = 'right'
  @Prop({ reflectToAttr: true }) aligned = 'left'
  @Prop()
  @Prop()
  backNav: boolean = true
  @Prop()
  btnProps: JSXElements.BearerButtonAttributes = {}

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
      <bearer-button {...this.btnProps} kind={this.kind} onClick={this.toggleDisplay}>
        {this.kind}
        <slot name="btn-content" />
      </bearer-button>,

      <div class={`popover direction-${this.direction} ${!this.visible && 'hidden'}`}>
        <h3 class="popover-header">
          {this.backNav && <bearer-navigator-back class="header-arrow" />}
          <span class="header">{this.header}</span>
        </h3>
        <div class="popover-body">
          <slot />
        </div>
      </div>
    ]
  }
}
