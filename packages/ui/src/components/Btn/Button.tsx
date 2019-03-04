/**
 * Temporary button for the lib.
 * No shadow on this one!
 */

import { Component, Prop } from '@bearer/core'

export type BKind = 'embed' | 'primary' | 'secondary' | 'error'

@Component({
  tag: 'bearer-btn',
  styleUrl: '../Button/Button.scss'
})
export class Btn {
  @Prop() content: any
  @Prop() type: string
  @Prop({ reflectToAttr: true }) kind: BKind = 'primary'
  @Prop() as: string = 'button'

  // @Prop() size: 'medium' | 'small' | 'large' = 'medium'
  // @Prop() outline: boolean = false
  @Prop() disabled: boolean = false

  @Prop({ reflectToAttr: true }) dataTooltip: string
  @Prop({ reflectToAttr: true }) dataTooltipType: string

  render() {
    // tslint:disable-next-line variable-name
    const Tag = this.as as any
    // const classes = ['btn', `btn${this.outline ? '-outline' : ''}-${this.kind}`, `btn-${this.size}`]
    const classes = ['btn', `btn-${this.kind}`]

    return (
      <Tag
        class={classes.join(' ')}
        disabled={this.disabled}
        type={this.type}
        data-tooltip={this.dataTooltip}
        data-tooltip-type={this.dataTooltipType}
      >
        {this.content || <slot />}
      </Tag>
    )
  }
}
