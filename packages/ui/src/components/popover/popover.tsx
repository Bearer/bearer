import { Component, Event, EventEmitter, Method, Prop, State, Watch } from '@bearer/core'

import debug from '../../logger'
const logger = debug('bearer-popover')

export type TAlignement = 'left' | 'right'
export type TDirection = 'left' | 'right' | 'top' | 'bottom'

@Component({
  tag: 'bearer-popover',
  styleUrl: 'popover.scss',
  shadow: true
})
export class BearerPopover {
  @Event()
  visibilityChange: EventEmitter

  @Prop({ reflectToAttr: true, mutable: true }) opened: boolean = false
  @Prop({ reflectToAttr: true }) direction: TDirection = 'right'
  @Prop({ reflectToAttr: true }) aligned: TAlignement = 'left'

  @Prop() heading: string
  @Prop() content: string

  @Prop() backNav: boolean = false

  @State() _visible: boolean = false

  set visible(newValue: boolean) {
    if (newValue !== null && this._visible !== newValue) {
      logger('Popover: visibilityChangeHandler %s', newValue)
      this._visible = newValue
      this.visibilityChange.emit({ visible: this._visible })
    }
  }

  get visible(): boolean {
    return typeof this._visible === 'boolean' ? this._visible : false
  }

  clickInsidePopover(ev) {
    ev.stopImmediatePropagation()
  }

  toggleDisplay = e => {
    e.preventDefault()
    this.toggle()
  }

  @Method()
  toggle(newValue?: boolean) {
    // Set visibility to toggle param
    // or inverse the current one.

    this.visible = typeof newValue !== 'undefined' ? newValue : !this.visible
    this.opened = this.visible
  }

  @Watch('opened')
  watchOpened(newValue: boolean) {
    if (newValue === this.visible) {
      return
    }
    if (newValue === null || newValue === undefined) {
      this.opened = false
    }

    this.toggle(newValue)
  }

  componentDidLoad() {
    this._visible = this.opened
  }

  render() {
    return [
      // @ts-ignore
      <slot name="popover-toggler" onClick={this.toggleDisplay}>
        <bearer-button kind="embed">
          <slot name="popover-button" />
        </bearer-button>
      </slot>,

      <div
        onClick={this.clickInsidePopover}
        class={`popover direction-${this.direction} ${!this.visible && 'hidden'} aligned-${this.aligned}`}
      >
        <slot name="popover-container">
          <div class="popover-container">
            <div class="popover-header">
              <slot name="popover-header">
                {this.backNav && <bearer-navigator-back class="header-arrow" />}
                {this.heading && <h3>{this.heading}</h3>}
              </slot>
            </div>
            <div class="popover-content">
              <slot>{this.content}</slot>
              <div class="popover-actions">
                <slot name="popover-actions" />
              </div>
            </div>
            <slot name="popover-footer">
              <div class="popover-copyright">
                Powered by{' '}
                <a href="https://bearer.sh" target="_blank">
                  Bearer.sh
                </a>
              </div>
            </slot>
          </div>
        </slot>
      </div>
    ]
  }
}
