import { Component, Prop } from '@bearer/core'

@Component({
  tag: 'bearer-badge',
  styleUrl: './Badge.scss',
  shadow: true
})
export class BearerBadge {
  @Prop() content: any
  @Prop() kind: 'primary' | 'secondary' | 'success' | 'danger' | 'warning' | 'info' | 'light' | 'dark'

  render() {
    return <span class={`badge badge-${this.kind}`}>{this.content || <slot />}</span>
  }
}
