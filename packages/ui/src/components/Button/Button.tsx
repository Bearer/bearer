import { Component, Prop } from '@bearer/core'

@Component({
  tag: 'bearer-button',
  styleUrl: 'Button.scss',
  shadow: true
})
export class Button {
  @Prop()
  kind:
    | 'primary'
    | 'secondary'
    | 'success'
    | 'danger'
    | 'warning'
    | 'info'
    | 'light'
    | 'dark'
  @Prop() size: 'md' | 'sm' | 'lg' = 'md'
  @Prop() as: string = 'button'

  render() {
    const Tag = this.as
    const classes = ['btn', `btn-${this.kind}`, `btn-${this.size}`]
    return (
      <Tag class={classes.join(' ')}>
        <slot />
      </Tag>
    )
  }
}
