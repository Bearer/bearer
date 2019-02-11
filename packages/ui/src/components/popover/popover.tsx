import { Component, Event, EventEmitter, Method, Prop, State, Watch } from '@bearer/core'

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

  @Prop() button: string
  @Prop() title: string
  @Prop() content: string

  @Prop() backNav: boolean = false

  @State() _visible: boolean = false

  toggleDisplay = e => {
    e.preventDefault()

    console.debug('[BEARER]', 'Popover: toggleDisplay', !this.visible)
    this.visible = !this.visible
  }

  set visible(newValue: boolean) {
    if (newValue !== null && this._visible !== newValue) {
      console.debug('[BEARER]', 'Popover: visibilityChangeHandler', newValue)
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
  toggle(newValue?: boolean) {
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
      // @ts-ignore
      <slot name="popover-toggler" onClick={this.toggleDisplay}>
        <bearer-button kind="primary">{this.button}</bearer-button>
      </slot>,

      <div
        onClick={this.clickInsidePopover}
        class={`popover direction-${this.direction} ${!this.visible && 'hidden'} aligned-${this.aligned}`}
      >
        <slot>
          <div class="popover-container">
            <div class="popover-header">
              <slot name="popover-header">
                {this.backNav && <bearer-navigator-back class="header-arrow" />}
                <div class="popover-title">
                  <slot name="popover-title">
                    <h3>{this.title}</h3>
                  </slot>
                </div>
              </slot>
            </div>
            <div class="popover-content">
              <slot name="popover-content">{this.content}</slot>
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
