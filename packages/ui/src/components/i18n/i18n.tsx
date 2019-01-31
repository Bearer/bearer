import { Component, t, Prop } from '@bearer/core'

@Component({
  tag: 'bearer-i18n'
})
export class BearerI18n {
  @Prop() key: string
  @Prop() default?: string
  render() {
    return t(this.key, this.default)
  }
}
