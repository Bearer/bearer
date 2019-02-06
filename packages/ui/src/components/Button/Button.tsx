import { Component, Prop } from '@bearer/core'

export type BKind = 'action' | 'primary' | 'secondary' | 'success' | 'danger' | 'warning' | 'info' | 'light' | 'dark'

@Component({
  tag: 'bearer-button',
  styleUrl: 'Button.scss',
  shadow: true
})
export class Button {
  @Prop() content: any
  @Prop() type: string
  @Prop({ reflectToAttr: true }) kind: BKind = 'primary'
  @Prop() size: 'medium' | 'small' | 'large' = 'medium'
  @Prop() as: string = 'button'
  @Prop() outline: boolean = false
  @Prop() disabled: boolean = false
  @Prop({ reflectToAttr: true }) dataTooltip: string
  @Prop({ reflectToAttr: true }) dataTooltipType: string

  render() {
    // tslint:disable-next-line variable-name
    const Tag = this.as as any
    const classes = ['btn', `btn${this.outline ? '-outline' : ''}-${this.kind}`, `btn-${this.size}`]
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
