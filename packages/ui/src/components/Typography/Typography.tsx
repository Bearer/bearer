import { Component, Prop } from '@bearer/core'

@Component({
  tag: 'bearer-typography',
  styleUrl: 'Typography.scss',
  shadow: true
})
export class Typography {
  @Prop() as: string = 'p'
  @Prop()
  kind:
    | ''
    | 'h1'
    | 'h2'
    | 'h3'
    | 'h4'
    | 'h5'
    | 'h6'
    | 'text-muted'
    | 'display-1'
    | 'display-2'
    | 'display-3'
    | 'display-4' =
    ''

  render() {
    const Tag = this.as
    return (
      <Tag class={this.kind}>
        <slot />
      </Tag>
    )
  }
}
